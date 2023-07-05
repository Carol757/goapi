package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
)

func TestGetID(t *testing.T) {
	payload := []byte(`{
		"retailer": "ABC Store",
		"purchaseDate": "2023-01-01",
		"purchaseTime": "14:30",
		"items": [
			{
				"shortDescription": "Item 1",
				"price": "10.00"
			},
			{
				"shortDescription": "Item 2",
				"price": "20.00"
			}
		],
		"total": "30.00"
	}`)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", getID)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", rr.Code)
	}

	var id ID
	err = json.Unmarshal(rr.Body.Bytes(), &id)
	if err != nil {
		t.Fatal(err)
	}

	if id.ID == "" {
		t.Error("Empty ID")
	}
}

func TestGetPointsByID(t *testing.T) {

	id := "e9c206d2c42b620a8fd7ad254cb0c204c1278413722e27b37f477438f347d924"

	req, err := http.NewRequest("GET", "/receipts/"+id+"/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/receipts/{id}/points", getPointsByID)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", rr.Code)
	}

	var points Points
	err = json.Unmarshal(rr.Body.Bytes(), &points)
	if err != nil {
		t.Fatal(err)
	}

	if points.Points == -1 {
		t.Error("Negative points")
	}
}
