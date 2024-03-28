package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

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
