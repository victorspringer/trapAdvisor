package persistence

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/trip"
)

func TestTripRepository(t *testing.T) {
	var err error
	database.DB, err = database.Init("test")
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	defer func() {
		_, err = database.DB.Exec("DROP DATABASE IF EXISTS trapAdvisor_test")
		if err != nil {
			t.Errorf("%v", err)
			t.FailNow()
		}
		database.DB.Close()
	}()

	type args struct {
		t           *trip.Trip
		tripID      int
		travellerID int
	}
	tests := []struct {
		name     string
		r        *tripRepository
		args     args
		wantName string
		wantErr  bool
	}{
		{
			"ValidTripInsertion",
			&tripRepository{},
			args{
				&trip.Trip{
					ID:          20,
					Name:        "Test",
					StartDate:   "2007-02-02 02:02:20",
					EndDate:     "2007-12-02 02:02:20",
					Rating:      5.5,
					Review:      "Test",
					TravellerID: 123,
				},
				20,
				123,
			},
			"Test",
			false,
		},
		{
			"ValidTripInsertion",
			&tripRepository{},
			args{
				&trip.Trip{
					Name:        "Test2",
					StartDate:   "2007-02-02 02:02:20",
					EndDate:     "2007-12-02 02:02:20",
					Rating:      5.5,
					Review:      "Test",
					TravellerID: 123,
				},
				21,
				123,
			},
			"Test2",
			false,
		},
		{
			"ValidTripUpdate",
			&tripRepository{},
			args{
				&trip.Trip{
					ID:          1,
					Name:        "Updated",
					StartDate:   "2007-02-02 02:02:20",
					EndDate:     "2007-12-02 02:02:20",
					Rating:      5.5,
					Review:      "Test",
					TravellerID: 123,
				},
				1,
				123,
			},
			"Updated",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = tt.r.Store(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("tripRepository.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			trip, err := tt.r.Find(tt.args.tripID)
			if (err != nil) != tt.wantErr {
				t.Errorf("tripRepository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if trip.Name != tt.wantName {
				t.Errorf("invalid result")
				return
			}
			trips, err := tt.r.FindByTravellerID(tt.args.travellerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("tripRepository.FindByTravellerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tripIDInSlice(tt.args.tripID, trips) {
				t.Errorf("invalid result")
				return
			}
		})
	}
}

func tripIDInSlice(id int, list []*trip.Trip) bool {
	for _, b := range list {
		if b.ID == id {
			return true
		}
	}
	return false
}
