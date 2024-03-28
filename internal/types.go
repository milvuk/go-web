package internal

import "time"

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type Product struct {
	ID                 int64     `json:"id,string"`
	ProductName        string    `json:"productName"`
	ProductDescription string    `json:"productDescription"`
	Department         string    `json:"department"`
	Price              float32   `json:"price,string"`
	CreatedAt          time.Time `json:"createdAt"`
}
