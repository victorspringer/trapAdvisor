package friendship

// Repository provides access to a Friendship store.
type Repository interface {
	Store(*Friendship) error
	FindByTravellerID(int) ([]*Friendship, error)
}

// Friendship class.
type Friendship struct {
	TravellerID int `json:"travellerId"`
	FriendID    int `json:"friendId"`
	// Similarity  float64 `json:"similarity"`
}
