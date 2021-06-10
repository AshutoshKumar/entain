package service

import (
	"golang.org/x/net/context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	//update status based on the adversisement time
	for _, race := range races {
		today := time.Now()
		ts, _ := ptypes.Timestamp(race.AdvertisedStartTime)
		after := ts.After(today)
		if after {
			race.Status = "OPEN"
		}else {
			race.Status = "CLOSED"
		}
	}

	return &racing.ListRacesResponse{Races: races}, nil
}
