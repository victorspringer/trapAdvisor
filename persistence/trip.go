package persistence

import (
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/trip"
)

type tripRepository struct{}

func (r *tripRepository) Store(t *trip.Trip) error {
	if t.ID != 0 {
		stmt, err := database.DB.Prepare("SELECT id FROM trip WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		var id int
		if err = stmt.QueryRow(t.ID).Scan(&id); err != nil {
			insert, err := database.DB.Prepare("INSERT INTO trip VALUES( ?, ?, ?, ?, ?, ?, ? )")
			if err != nil {
				return err
			}
			defer insert.Close()

			_, err = insert.Exec(t.ID, t.Name, t.StartDate, t.EndDate, t.Rating, t.Review, t.TravellerID)
			if err != nil {
				return err
			}
		} else {
			update, err := database.DB.Prepare(`
				UPDATE trip SET
					id = ?,
					name = ?,
					startDate = ?,
					endDate = ?,
					rating = ?,
					review = ?,
					travellerId = ?
				WHERE id = ?
			`)
			if err != nil {
				return err
			}
			defer update.Close()

			_, err = update.Exec(t.ID, t.Name, t.StartDate, t.EndDate, t.Rating, t.Review, t.TravellerID, t.ID)
			if err != nil {
				return err
			}
		}
	} else {
		insert, err := database.DB.Prepare("INSERT INTO trip VALUES( ?, ?, ?, ?, ?, ?, ? )")
		if err != nil {
			return err
		}
		defer insert.Close()

		_, err = insert.Exec(nil, t.Name, t.StartDate, t.EndDate, t.Rating, t.Review, t.TravellerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *tripRepository) Find(id int) (*trip.Trip, error) {
	stmt, err := database.DB.Prepare("SELECT * FROM trip WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var t trip.Trip
	if err = stmt.QueryRow(id).Scan(
		&t.ID, &t.Name, &t.StartDate, &t.EndDate, &t.Rating, &t.Review, &t.TravellerID,
	); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *tripRepository) FindByTravellerID(travellerID int) ([]*trip.Trip, error) {
	stmt, err := database.DB.Prepare("SELECT * FROM trip WHERE travellerId = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(travellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []*trip.Trip
	for rows.Next() {
		var t trip.Trip
		if err = rows.Scan(
			&t.ID, &t.Name, &t.StartDate, &t.EndDate, &t.Rating, &t.Review, &t.TravellerID,
		); err != nil {
			return nil, err
		}
		trips = append(trips, &t)
	}

	return trips, nil
}

// NewTripRepository returns a new instance of a MySQL trip repository.
func NewTripRepository() trip.Repository {
	return &tripRepository{}
}
