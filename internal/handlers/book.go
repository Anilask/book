package handlers

import (
	"encoding/json"
	"library/internal/utils"
	"library/pkg/jwt"
	"net/http"
	"strconv"
	"time"
)

// Book represents the structure for a book entry
type Book struct {
	Name            string `json:"name"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publicationYear"`
}

// AddBook adds a new book to the library, accessible only by admin users
func AddBook(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.ExtractClaims(r)
	if err != nil {
		http.Error(w, "Unauthorized - Token required", http.StatusUnauthorized)
		return
	}

	// Check if the user is an admin
	if claims.UserType != "admin" {
		http.Error(w, "Unauthorized - Access restricted to admins", http.StatusForbidden)
		return
	}

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the book details
	if book.Name == "" || book.Author == "" || book.PublicationYear < 1000 || book.PublicationYear > time.Now().Year() {
		http.Error(w, "Invalid book details", http.StatusBadRequest)
		return
	}

	// Prepare the record for CSV
	record := []string{book.Name, book.Author, strconv.Itoa(book.PublicationYear)}
	if err := utils.AppendToCSV("/data/ezyzip/library-management/data/regularUser.csv", record); err != nil {
		http.Error(w, "Failed to add the book to the library", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// DeleteBook removes a book from the library, accessible only by admin users
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.ExtractClaims(r)
	if err != nil {
		http.Error(w, "Unauthorized - Token required", http.StatusUnauthorized)
		return
	}

	// Check if the user is an admin
	if claims.UserType != "admin" {
		http.Error(w, "Unauthorized - Access restricted to admins", http.StatusForbidden)
		return
	}

	var book struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the book name
	if book.Name == "" {
		http.Error(w, "Invalid book name", http.StatusBadRequest)
		return
	}

	// Attempt to remove the book from the CSV
	if err := utils.RemoveFromCSV("/data/ezyzip/library-management/data/regularUser.csv", book.Name); err != nil {
		http.Error(w, "Failed to delete the book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted successfully"})
}
