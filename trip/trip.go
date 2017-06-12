package trip

// Trip class.
type Trip struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	StartDate   string  `json:"startDate"`
	EndDate     string  `json:"endDate"`
	Rating      float64 `json:"rating"`
	Review      string  `json:"review"`
	TravellerID int     `json:"travellerId"`
}

// Repository provides access to a Trip store.
type Repository interface {
	Store(*Trip) error
	Find(int) (*Trip, error)
	FindByTravellerID(int) ([]*Trip, error)
}
