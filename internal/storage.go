package internal

import (
	"database/sql"
)

func albums(db *sql.DB) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		var p sql.NullFloat64
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &p); err != nil {
			return nil, err
		}
		if p.Valid {
			alb.Price = float32(p.Float64)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return albums, nil
}

func albumByID(db *sql.DB, id int64) (Album, error) {
	var alb Album
	var p sql.NullFloat64

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &p); err != nil {
		return alb, err
	}
	if p.Valid {
		alb.Price = float32(p.Float64)
	}
	return alb, nil
}

func addAlbum(db *sql.DB, alb Album) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func deleteAlbum(db *sql.DB, id int64) error {
	res, err := db.Exec("DELETE FROM album WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func updateAlbum(db *sql.DB, id int64, alb Album) error {
	res, err := db.Exec("UPDATE album SET title=$1, artist=$2, price=$3 WHERE id=$4",
		alb.Title, alb.Artist, alb.Price, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
