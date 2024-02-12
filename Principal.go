package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// principal Table
type Principal struct {
	gorm.Model
	// Principal_Id  string `json:"principalid" gorm:"primaryKey"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	Email         string `json:"email"`
	Qualification string `json:"qualification"`
}

// Get principals
func GetPrincipals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var principal []Principal
	DB.Find(&principal)
	json.NewEncoder(w).Encode(principal)
}

// Get principal
func GetPrincipal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var principal Principal
	params := mux.Vars(r)
	DB.First(&principal, params["principalid"])
	json.NewEncoder(w).Encode(principal)
}

// Add a principal into principal table
func AddPrincipal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var principal Principal
	json.NewDecoder(r.Body).Decode(&principal)
	DB.Create(&principal)
	json.NewEncoder(w).Encode(principal)
}
