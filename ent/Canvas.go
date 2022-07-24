package ent

import (
	"database/sql"
	"log"
)

type Canvas struct {
	ID           string `json:"id"`
	Drawing      string `json:"drawing"`
	CreationDate string `json:"creationDate"`
}

func (canvas *Canvas) GetCanvas(db *sql.DB) error {
	return db.QueryRow("SELECT drawing, creationDate FROM canvas WHERE id=?",
		canvas.ID).Scan(&canvas.Drawing, &canvas.CreationDate)
}

func (canvas *Canvas) UpdateCanvas(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE canvas SET drawing=?, creationDate=? WHERE id=?",
			canvas.Drawing, canvas.CreationDate, canvas.ID)

	return err
}

func (canvas *Canvas) DeleteCanvas(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM canvas WHERE id=?", canvas.ID)

	return err
}

func (canvas *Canvas) CreateCanvas(db *sql.DB) error {

	_, err := db.Exec("INSERT INTO canvas (id,drawing, creationdate) VALUES (?, ?,?)", canvas.ID, canvas.Drawing,canvas.CreationDate)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetCanvasList(db *sql.DB) ([]Canvas, error) {
	rows, err := db.Query(
		"SELECT id, drawing, creationdate FROM canvas")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	canvasList := []Canvas{}

	for rows.Next() {
		var canvas Canvas
		if err := rows.Scan(&canvas.ID, &canvas.Drawing, &canvas.CreationDate); err != nil {
			return nil, err
		}
		canvasList = append(canvasList, canvas)
	}

	return canvasList, nil
}
