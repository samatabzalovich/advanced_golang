package main

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"product-service/data"
	"strconv"
	"time"
)

func (app *Config) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	// Retrieve the value of the 'minPrice' parameter
	minPriceStr := queryParams.Get("minPrice")

	// Retrieve the value of the 'maxPrice' parameter
	maxPriceStr := queryParams.Get("maxPrice")
	var minPrice float64
	var err error
	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid 'minPrice' parameter", http.StatusBadRequest)
			return
		}
	}

	// Process the 'maxPrice' parameter
	var maxPrice float64
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid 'maxPrice' parameter", http.StatusBadRequest)
			return
		}
	}

	products, err := app.Models.LogEntry.Get(maxPrice, minPrice)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
	err = app.writeJSON(w, http.StatusOK, products)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}

}

func (app Config) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Name  string  `json:"name"`
		Price float32 `json:"price"`
		Data  string  `json:"data"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	product := data.Product{Name: requestPayload.Name, Price: requestPayload.Price, Data: requestPayload.Data, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	err = app.Models.LogEntry.Insert(product)
	if err != nil {
		app.errorJSON(w, errors.New("cerver err"), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("createg product %s", product.Name),
		Data:    product,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
