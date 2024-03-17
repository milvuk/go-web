package main

import (
	"database/sql"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func albums(db *sql.DB) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, err
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return albums, nil
}

func albumsByArtist(db *sql.DB, name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, err
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

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return alb, err
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
