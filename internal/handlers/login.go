package handlers

import (
	"encoding/json"
	"library/internal/model"
	"library/pkg/jwt"
	"net/http"
)

// Static list of users for demonstration purposes
var users = []model.User{
	{Username: "admin", UserType: "admin", Password: "admin123"},
	{Username: "user", UserType: "regular", Password: "user123"},
}

// Login processes the login request and returns a JWT token if successful
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials model.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	for _, user := range users {
		if user.Username == credentials.Username && user.Password == credentials.Password {
			tokenString, err := jwt.GenerateToken(user.Username, user.UserType)
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
			return
		}
	}

	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
}
