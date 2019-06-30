package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Location struct {
	Latitude	float64 `json:"lat"`
	Longitude	float64 `json:"lng"`
}

type Country struct {
	CountryCode	string `json:"country_code"` // ISO 3166-1 alpha-2 2-letter country code
}

func georeverseHandler(w http.ResponseWriter, r *http.Request) {
	var latlng Location
	var results Country

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&latlng)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	results.CountryCode = reverser.GetCountryCode(latlng.Longitude, latlng.Latitude)

	if results.CountryCode == "" {
		log.Printf("Unable to decode country for %+v\n", latlng)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
