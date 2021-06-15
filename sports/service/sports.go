package service

import (
	"time"
	"golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes"
	"github.com/AshutoshKumar/entain/sports/db"
	"github.com/AshutoshKumar/entain/sports/proto/sports"
)

type Sports interface {
	// ListRaces will return a collection of races.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)

}

// racingService implements the Racing interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

func (s sportsService) ListEvents(_ context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	sportsList, err := s.sportsRepo.ListEvents(in.Filter)
	if err != nil {
		return nil, err
	}

	//update status based on the adversisement time
	for _, sport := range sportsList {
		today := time.Now()
		ts, _ := ptypes.Timestamp(sport.AdvertisedStartTime)
		after := ts.After(today)
		if after {
			sport.Status = "OPEN"
		} else {
			sport.Status = "CLOSED"
		}
	}

	return &sports.ListEventsResponse{Sports: sportsList}, nil
}

// NewRacingService instantiates and returns a new racingService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}
