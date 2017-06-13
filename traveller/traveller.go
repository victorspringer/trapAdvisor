package traveller

// Traveller class.
type Traveller struct {
	Name string `json:"name"`
	ID   int    `json:"id,string"`
}

// Repository provides access to a Traveller store.
type Repository interface {
	Store(*Traveller) error
	Find(int) (*Traveller, error)
}
