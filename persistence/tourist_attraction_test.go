package persistence

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/victorspringer/trapAdvisor/database"
	ta "github.com/victorspringer/trapAdvisor/touristattraction"
)

func TestTouristAttractionRepository(t *testing.T) {
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
		t                   *ta.TouristAttraction
		touristAttractionID int
		tripID              int
		namePart            string
	}
	tt := struct {
		name         string
		r            *touristAttractionRepository
		args         args
		wantName     string
		wantNamePart string
		wantErr      bool
	}{
		"ValidTouristAttractionInsertion",
		&touristAttractionRepository{},
		args{
			&ta.TouristAttraction{
				Name:      "Test",
				Location:  "Brazil",
				VisitDate: "2007-02-02 02:02:20",
				Rating:    5.5,
				Pros:      "Test",
				Cons:      "Test",
				TripID:    1,
			},
			2,
			1,
			"Eif",
		},
		"Test",
		"Torre Eiffel",
		false,
	}
	t.Run(tt.name, func(t *testing.T) {
		if err = tt.r.Store(tt.args.t); (err != nil) != tt.wantErr {
			t.Errorf("touristAttractionRepository.Store() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		ta, err := tt.r.Find(tt.args.touristAttractionID)
		if (err != nil) != tt.wantErr {
			t.Errorf("touristAttractionRepository.Find() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if ta.Name != tt.wantName {
			t.Errorf("invalid result")
			return
		}
		tas, err := tt.r.FindByTripID(tt.args.tripID)
		if (err != nil) != tt.wantErr {
			t.Errorf("touristAttractionRepository.FindByTripID() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !touristAttractionIDInSlice(tt.args.touristAttractionID, tas) {
			t.Errorf("invalid result")
			return
		}
		tasNamePart, err := tt.r.FindByNamePart(tt.args.namePart)
		if (err != nil) != tt.wantErr {
			t.Errorf("touristAttractionRepository.FindByNamePart() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if tasNamePart[0][0] != tt.wantNamePart {
			t.Errorf("invalid result")
			return
		}
		tasMostVisited, err := tt.r.FindMostVisited()
		if (err != nil) != tt.wantErr {
			t.Errorf("touristAttractionRepository.FindMostVisited() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if tasMostVisited[0][0] != tt.wantName {
			fmt.Println(tasMostVisited[0][0], tt.wantName)
			t.Errorf("invalid result")
			return
		}
		tasBestRated, err := tt.r.FindBestRated()
		if (err != nil) != tt.wantErr {
			t.Errorf("touristAttractionRepository.FindBestRated() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if tasBestRated[1][0] != tt.wantName {
			t.Errorf("invalid result")
			return
		}
	})
}

func touristAttractionIDInSlice(id int, list []*ta.TouristAttraction) bool {
	for _, b := range list {
		if b.ID == id {
			return true
		}
	}
	return false
}
