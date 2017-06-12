package persistence

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/friendship"
)

func TestFriendshipRepository(t *testing.T) {
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
		t *friendship.Friendship
	}
	tests := []struct {
		name    string
		r       *friendshipRepository
		args    args
		want    int
		wantErr bool
	}{
		{
			"ValidFriendshipInsertion",
			&friendshipRepository{},
			args{
				&friendship.Friendship{TravellerID: 123, FriendID: 1234},
			},
			123,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = tt.r.Store(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("friendshipRepository.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			friendships, err := tt.r.FindByTravellerID(1234)
			if (err != nil) != tt.wantErr {
				t.Errorf("friendshipRepository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !friendIDInSlice(tt.want, friendships) {
				t.Errorf("invalid result")
				return
			}
		})
	}
}

func friendIDInSlice(id int, list []*friendship.Friendship) bool {
	for _, b := range list {
		if b.FriendID == id {
			return true
		}
	}
	return false
}
