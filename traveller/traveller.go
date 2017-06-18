package traveller

// Traveller class.
type Traveller struct {
	ID           int    `json:"id,string"`
	Name         string `json:"name"`
	SessionToken string `json:"-"`
}

// Repository provides access to a Traveller store.
type Repository interface {
	Store(*Traveller) error
	Find(int) (*Traveller, error)
	FindBySessionToken(int, string) error
}
