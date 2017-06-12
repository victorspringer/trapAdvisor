package touristattraction

// TouristAttraction class.
type TouristAttraction struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	VisitDate string  `json:"visitDate"`
	Rating    float64 `json:"rating"`
	Pros      string  `json:"pros"`
	Cons      string  `json:"cons"`
	TripID    int     `json:"tripId"`
}

// Repository provides access to a TouristAttraction store.
type Repository interface {
	Store(*TouristAttraction) error
	Find(int) (*TouristAttraction, error)
	FindByTripID(int) ([]*TouristAttraction, error)
	FindByNamePart(string) (map[int][2]string, error)
	FindMostVisited() (map[int][3]string, error)
	FindBestRated() (map[int][3]string, error)
}
