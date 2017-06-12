package persistence

import (
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/traveller"
)

type travellerRepository struct{}

func (r *travellerRepository) Store(t *traveller.Traveller) error {
	stmt, err := database.DB.Prepare("SELECT id FROM traveller WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int
	if err = stmt.QueryRow(t.ID).Scan(&id); err != nil {
		insert, err := database.DB.Prepare("INSERT INTO traveller VALUES( ?, ? )")
		if err != nil {
			return err
		}
		defer insert.Close()

		_, err = insert.Exec(t.ID, t.Name)
		if err != nil {
			return err
		}
	} else {
		update, err := database.DB.Prepare("UPDATE traveller SET id = ?, name = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer update.Close()

		_, err = update.Exec(id, t.Name, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *travellerRepository) Find(id int) (*traveller.Traveller, error) {
	stmt, err := database.DB.Prepare("SELECT * FROM traveller WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var t traveller.Traveller
	if err = stmt.QueryRow(id).Scan(&t.ID, &t.Name); err != nil {
		return nil, err
	}

	return &t, nil
}

// NewTravellerRepository returns a new instance of a MySQL traveller repository.
func NewTravellerRepository() traveller.Repository {
	return &travellerRepository{}
}
