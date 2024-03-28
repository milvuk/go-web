package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Product struct {
	ID                 int64     `json:"id,string"`
	ProductName        string    `json:"productName"`
	ProductDescription string    `json:"productDescription"`
	Department         string    `json:"department"`
	Price              float32   `json:"price,string"`
	CreatedAt          time.Time `json:"createdAt"`
}

var ErrNotFound = errors.New("not found")

func retrieveProducts() ([]Product, error) {
	url := mockApiAddress() + "/products"

	products := []Product{}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, err
}

func retrieveProduct(id int64) (Product, error) {
	url := mockApiAddress() + "/products/" + strconv.FormatInt(id, 10)

	var product Product

	resp, err := http.Get(url)
	if err != nil {
		return product, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return product, ErrNotFound
	}

	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

func mockApiAddress() string {
	return "https://" + ViperEnvVariable("MOCKAPI_PROJECT_KEY") + ".mockapi.io/api/v1"
}
