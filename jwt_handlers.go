package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = []byte("secret_key")

// Used to check credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Login Page
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	params := mux.Vars(r)
	username := params["username"]
	password := params["password"]
	HashedPassword, err := GenerateHash(password)
	if err != nil {
		http.Error(w, "Could not generate password while Adding Credentials", http.StatusForbidden)
		return
	}
	result := DB.Where("username = ? AND password = ?", username, HashedPassword).First(&credentials)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Credential does not exist", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))
}

// This credentials will be added by Principle while adding new students
func AddCredentials(username, password, role string) (Credentials, error) {
	var credentials Credentials
	HashedPassword, err := GenerateHash(password)
	if err != nil {
		log.Printf("Error while generating credentials: %v", err)
		return credentials, err
	}
	// role := result["role"]
	credentials.Username = username
	credentials.Password = HashedPassword
	credentials.Role = role
	errCreate := DB.Create(&credentials)
	if errCreate.Error != nil {
		log.Printf("Error creating credentials: %v", errCreate)
		return credentials, err
	}

	// json.NewEncoder(w).Encode(credentials)
	return credentials, nil
}

// Function to generate Hash for Passwords
func GenerateHash(password string) (string, error) {
	// Convert password string to byte slice
	passwordBytes := []byte(password)

	// Generate hash of the password
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert hashed bytes to string and return
	hashedPassword := string(hashedBytes)
	return hashedPassword, nil
}
