package main

import (
	"database/sql"
	"github.com/lib/pq"
)

type shoe struct {
	ID                    int     `json:"id"`
	Name                  string  `json:"shoe"`
	TrueToSizeData        []int64 `json:"trueToSizeData"`
	TrueToSizeCalculation float64 `json:"trueToSizeCalculation"`
}

func getShoes(db *sql.DB) ([]shoe, error) {
	rows, err := db.Query(
		"SELECT id, name, \"trueToSizeData\", \"trueToSizeCalculation\" FROM shoes")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shoes := []shoe{}

	for rows.Next() {
		var s shoe
		if err := rows.Scan(&s.ID, &s.Name, pq.Array(&s.TrueToSizeData), &s.TrueToSizeCalculation); err != nil {
			return nil, err
		}
		shoes = append(shoes, s)
	}

	return shoes, nil
}

func getOrCreateShoeByName(db *sql.DB, name string) (shoe, error) {
	var s shoe
	err := db.QueryRow(
		"SELECT id, name, \"trueToSizeData\", \"trueToSizeCalculation\" FROM shoes WHERE LOWER(name)=LOWER($1)",
		name,
	).Scan(&s.ID, &s.Name, pq.Array(&s.TrueToSizeData), &s.TrueToSizeCalculation)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			s, err = createShoe(db, name)
		default:
			return s, err
		}
	}

	return s, nil
}

func (s *shoe) findById(db *sql.DB) error {
	return db.QueryRow(
		"SELECT * FROM shoes WHERE id=$1",
		&s.ID,
	).Scan(&s.ID, &s.Name, pq.Array(&s.TrueToSizeData), &s.TrueToSizeCalculation)
}

func createShoe(db *sql.DB, name string) (shoe, error) {
	var s shoe
	err := db.QueryRow(
		"INSERT INTO shoes(name) VALUES($1) RETURNING id, name, 'trueToSizeData', 'trueToSizeCalculation'",
		name).Scan(&s.ID, &s.Name, pq.Array(&s.TrueToSizeData), &s.TrueToSizeCalculation)

	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *shoe) updateShoe(db *sql.DB, newValue int) error {
	_, err :=
		db.Exec("UPDATE shoes SET \"trueToSizeData\" = array_append(\"trueToSizeData\", $1::int) WHERE id=$2",
			newValue, s.ID)

	// reload data
	s.findById(db)

	return err
}

func (s *shoe) deleteShoe(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM shoes WHERE id=$1", s.ID)

	return err
}
