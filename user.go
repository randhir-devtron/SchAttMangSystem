package main

import (
	// "encoding/json"
	// "fmt"
	// "net/http"

	// "github.com/gorilla/mux"
	// "gorm.io/driver/postgres"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

// var pid int = 0

const DNS = "user=postgres password=hello@1234 dbname=postgres sslmode=disable host=localhost port=5432"

// type User struct {
// 	gorm.Model
// 	FirstName string `json:"firstname"`
// 	LastName  string `json:"lastname"`
// 	Email     string `json:"email"`
// }

// InitialMigration Function to check if the Database is connecting or not
func InitialMigration() {
	DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	// DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Principal{})
	DB.AutoMigrate(&Teacher{})
	DB.AutoMigrate(&Student{})
	DB.AutoMigrate(&Teacher_Attendance{})
	DB.AutoMigrate(&Student_Attendance{})
	DB.AutoMigrate(&Credentials{})

}

// func GetUsers(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var users []User
// 	DB.Find(&users)
// 	json.NewEncoder(w).Encode(users)
// }

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var user User
// 	DB.First(&user, params["id"])
// 	json.NewEncoder(w).Encode(user)
// }

// func CreateUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var user User
// 	json.NewDecoder(r.Body).Decode(&user)
// 	DB.Create(&user)
// 	json.NewEncoder(w).Encode(user)
// }

// func UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var user User
// 	DB.First(&user, params["id"])
// 	json.NewDecoder(r.Body).Decode(&user)
// 	DB.Save(&user)
// 	json.NewEncoder(w).Encode(user)
// }

// func DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var user User
// 	DB.Delete(&user, params["id"])
// 	// json.NewDecoder(r.Body).Decode(&user)
// 	// DB.Save(&user)
// 	json.NewEncoder(w).Encode("The user is deleted successfully")
// }
