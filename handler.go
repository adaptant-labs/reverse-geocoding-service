package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

type Location struct {
	Latitude	float64 `json:"lat"`
	Longitude	float64 `json:"lng"`
}

type Country struct {
	CountryCode	string `json:"country_code" maxminddb:"iso_code"` // ISO 3166-1 alpha-2 2-letter country code
}

type MaxMindRecord struct {
	Country Country `maxminddb:"country"`
}

func georeverseIPHandler(w http.ResponseWriter, r *http.Request) {
	var record MaxMindRecord

	vars := mux.Vars(r)
	ip := net.ParseIP(vars["ipAddress"])
	if ip == nil {
		http.Error(w, "Invalid IP Address specified", http.StatusBadRequest)
		return
	}

	err := geodb.Lookup(ip, &record)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(record.Country)
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
