package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type pointsStruct struct {
	ID     uuid.UUID `json:"id"`
	Points int       `json:"points"`
}

var pointsHolder = []pointsStruct{}

func getPoints(w http.ResponseWriter, r *http.Request) {
	// use ParamsFromContext() to retrieve a slice containing the URL parameters in the request
	params := httprouter.ParamsFromContext(r.Context())

	// use ByName() to get the value of the id paramter from the slice. Convert it to uuid
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var points int

	// loop through pointsHolder and match pointsStruct to id
	for _, val := range pointsHolder {
		if val.ID == id {
			points = val.Points
			break
		} 
	}

	js := fmt.Sprintf(`{ "points": %d }`, points)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}

func processReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt struct {
		Retailer     string `json:"retailer"`
		PurchaseDate string `json:"purchaseDate"`
		PurchaseTime string `json:"purchaseTime"`
		Total        string `json:"total"`
		Items        []struct {
			ShortDescription string `json:"shortDescription"`
			Price            string `json:"price"`
		} `json:"items"`
	}

	// read from request body
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	points := 0

	// one point for every alphanumeric character in the retailer name
	points = countRetailerName(points, receipt.Retailer)

	// 50 points if the total is a round dollar amount with no cents
	points, err = roundDollarAmount(points, receipt.Total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 25 points if the total is a multiple of 0.25
	points, err = multipleOfPointTwoFive(points, receipt.Total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5 point for every two items on the receipt
	points = everyTwoItems(points, receipt.Items)

	// if the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
	// the result is the number of points earned
	points, err = trimmedTotal(points, receipt.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 6 points if the day in the purchase date is odd
	points, err = oddPurchaseDate(points, receipt.PurchaseDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	points, err = purchaseTime(points, receipt.PurchaseTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create uuid
	id := uuid.New()

	newPoints := pointsStruct{ID: id, Points: points}

	// store uuid and points in in-memory struct
	pointsHolder = append(pointsHolder, newPoints)

	js := fmt.Sprintf(`{ "id": %q }`, id)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}
