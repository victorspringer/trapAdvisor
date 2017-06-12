package persistence

import (
	"strconv"

	"github.com/victorspringer/trapAdvisor/database"
	ta "github.com/victorspringer/trapAdvisor/touristattraction"
)

type touristAttractionRepository struct{}

func (r *touristAttractionRepository) Store(t *ta.TouristAttraction) error {
	insert, err := database.DB.Prepare(
		"INSERT INTO touristAttraction VALUES( ?, ?, ?, ?, ?, ?, ?, ? )",
	)
	if err != nil {
		return err
	}
	defer insert.Close()

	_, err = insert.Exec(nil, t.Name, t.Location, t.VisitDate, t.Rating, t.Pros, t.Cons, t.TripID)
	if err != nil {
		return err
	}

	return nil
}

func (r *touristAttractionRepository) Find(id int) (*ta.TouristAttraction, error) {
	stmt, err := database.DB.Prepare("SELECT * FROM touristAttraction WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var t ta.TouristAttraction
	if err = stmt.QueryRow(id).Scan(
		&t.ID, &t.Name, &t.Location, &t.VisitDate, &t.Rating, &t.Pros, &t.Cons, &t.TripID,
	); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *touristAttractionRepository) FindByTripID(tripID int) ([]*ta.TouristAttraction, error) {
	stmt, err := database.DB.Prepare("SELECT * FROM touristAttraction WHERE tripId = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var touristAttractions []*ta.TouristAttraction
	for rows.Next() {
		var t ta.TouristAttraction
		if err = rows.Scan(
			&t.ID, &t.Name, &t.Location, &t.VisitDate, &t.Rating, &t.Pros, &t.Cons, &t.TripID,
		); err != nil {
			return nil, err
		}
		touristAttractions = append(touristAttractions, &t)
	}

	return touristAttractions, nil
}

func (r *touristAttractionRepository) FindByNamePart(namePart string) (map[int][2]string, error) {
	stmt, err := database.DB.Prepare(
		"SELECT name, location FROM touristAttraction WHERE name LIKE ? GROUP BY name, location",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	namePart = "%" + namePart + "%"
	rows, err := stmt.Query(namePart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	touristAttractions := make(map[int][2]string)
	idx := 0
	for rows.Next() {
		var name, location string
		if err = rows.Scan(&name, &location); err != nil {
			return nil, err
		}
		touristAttractions[idx] = [2]string{name, location}
		idx++
	}

	return touristAttractions, nil
}

func (r *touristAttractionRepository) FindMostVisited() (map[int][3]string, error) {
	stmt, err := database.DB.Prepare(`
		SELECT name, location, COUNT(*) AS count
		FROM touristAttraction
		GROUP BY name, location
		ORDER BY count DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	touristAttractions := make(map[int][3]string)
	idx := 0
	for rows.Next() {
		var name, location string
		var count int
		if err = rows.Scan(&name, &location, &count); err != nil {
			return nil, err
		}
		touristAttractions[idx] = [3]string{name, location, strconv.Itoa(count)}
		idx++
	}

	return touristAttractions, nil
}

func (r *touristAttractionRepository) FindBestRated() (map[int][3]string, error) {
	stmt, err := database.DB.Prepare(`
		SELECT name, location, AVG(rating) as average
		FROM touristAttraction
		GROUP BY name, location 
		ORDER BY average DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	touristAttractions := make(map[int][3]string)
	idx := 0
	for rows.Next() {
		var name, location string
		var avg float64
		if err = rows.Scan(&name, &location, &avg); err != nil {
			return nil, err
		}
		touristAttractions[idx] = [3]string{name, location, strconv.FormatFloat(avg, 'f', -1, 64)}
		idx++
	}

	return touristAttractions, nil
}

// NewTouristAttractionRepository returns a new instance of a MySQL touristAttraction repository.
func NewTouristAttractionRepository() ta.Repository {
	return &touristAttractionRepository{}
}
