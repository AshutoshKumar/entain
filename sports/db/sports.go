package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/golang/protobuf/ptypes"
	"github.com/AshutoshKumar/entain/sports/proto/sports"
)

// SportsRepo provides repository access to races.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// ListEvents will return a list of sports.
	ListEvents(filter *sports.ListEventsRequestFilter) ([]*sports.Sport, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

func (s sportsRepo) Init() error {
	var err error
	s.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = s.seed()
	})

	return err
}

func (s sportsRepo) ListEvents(filter *sports.ListEventsRequestFilter) ([]*sports.Sport, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportsQueries()[sportsList]
	query, args = s.applyFilter(query, filter)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanRaces(rows)

}

// NewSportsRepo creates a new races repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

func (s *sportsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.Id) > 0 {
		clauses = append(clauses, "id IN("+strings.Repeat("?,", len(filter.Id)-1)+"?)")

		for _, id := range filter.Id {
			args = append(args, id)
		}
	}

	if len(clauses) != 0 {
		clauses = append(clauses, "Order by advertised_start_time")
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (s *sportsRepo) scanRaces(
	rows *sql.Rows,
) ([]*sports.Sport, error) {
	var sportList []*sports.Sport

	for rows.Next() {
		var sport sports.Sport
		var advertisedStart time.Time

		if err := rows.Scan(&sport.Id, &sport.Name, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			//"error converting date to timestamp proto")
		}
		if err != nil {
			return nil, err
		}

		sport.AdvertisedStartTime = ts

		sportList = append(sportList, &sport)
	}

	return sportList, nil
}
