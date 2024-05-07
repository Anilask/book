package handlers

import (
	"encoding/json"
	"net/http"
	"library/internal/utils"
	"library/pkg/jwt"
)

// Home displays a list of books based on the user type
func Home(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.ExtractClaims(r)
	if err != nil {
		http.Error(w, "Unauthorized - Token required", http.StatusUnauthorized)
		return
	}

	// Get books available to all regular users
	regularBooks, err := utils.ReadCSV("/data/ezyzip/library-management/data/regularUser.csv")
	if err != nil {
		http.Error(w, "Failed to read books for regular users", http.StatusInternalServerError)
		return
	}

	var books []string
	books = append(books, regularBooks...)

	// If the user is an admin, add the admin-specific books
	if claims.UserType == "admin" {
		adminBooks, err := utils.ReadCSV("/data/ezyzip/library-management/data/adminUser.csv")
		if err != nil {
			http.Error(w, "Failed to read books for admin users", http.StatusInternalServerError)
			return
		}
		books = append(books, adminBooks...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
