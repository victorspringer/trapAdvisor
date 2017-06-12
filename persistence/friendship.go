package persistence

import (
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/friendship"
)

type friendshipRepository struct{}

func (r *friendshipRepository) Store(f *friendship.Friendship) error {
	insert, err := database.DB.Prepare("INSERT INTO friendship VALUES( ?, ?, ? )")
	if err != nil {
		return err
	}
	defer insert.Close()

	_, err = insert.Exec(f.TravellerID, f.FriendID, nil)
	if err != nil {
		return err
	}

	_, err = insert.Exec(f.FriendID, f.TravellerID, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *friendshipRepository) FindByTravellerID(travellerID int) ([]*friendship.Friendship, error) {
	//stmt, err := database.DB.Prepare("SELECT * FROM friendship WHERE travellerId = ?")
	stmt, err := database.DB.Prepare("SELECT travellerId, friendId FROM friendship WHERE travellerId = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(travellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*friendship.Friendship
	for rows.Next() {
		var f friendship.Friendship
		//if err = rows.Scan(&f.TravellerID, &f.FriendID, &f.Similarity); err != nil {
		if err = rows.Scan(&f.TravellerID, &f.FriendID); err != nil {
			return nil, err
		}
		friendships = append(friendships, &f)
	}

	return friendships, nil
}

// NewFriendshipRepository returns a new instance of a MySQL friendship repository.
func NewFriendshipRepository() friendship.Repository {
	return &friendshipRepository{}
}
