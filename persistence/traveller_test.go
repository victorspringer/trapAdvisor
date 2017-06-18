package persistence

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/traveller"
)

func TestTravellerRepository(t *testing.T) {
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
		t *traveller.Traveller
	}
	tests := []struct {
		name    string
		r       *travellerRepository
		args    args
		want    string
		wantErr bool
	}{
		{
			"ValidTravellerInsertion",
			&travellerRepository{},
			args{
				&traveller.Traveller{ID: 12345, Name: "Rudolf", SessionToken: "test12345"},
			},
			"Rudolf",
			false,
		},
		{
			"ValidTravellerUpdate",
			&travellerRepository{},
			args{
				&traveller.Traveller{ID: 12345, Name: "Timo", SessionToken: "test12345"},
			},
			"Timo",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = tt.r.Store(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("travellerRepository.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			traveller, err := tt.r.Find(12345)
			if (err != nil) != tt.wantErr {
				t.Errorf("travellerRepository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = tt.r.FindBySessionToken(12345, "test12345")
			if (err != nil) != tt.wantErr {
				t.Errorf("travellerRepository.FindBySessionToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if traveller.Name != tt.want {
				t.Errorf("invalid result")
				return
			}
		})
	}
}
