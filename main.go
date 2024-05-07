package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "library/internal/auth"
    "library/internal/handlers"
)

func main() {
    r := mux.NewRouter()

    // Apply the ValidateTokenMiddleware
    r.HandleFunc("/login", handlers.Login).Methods("POST")
    r.HandleFunc("/home", auth.ValidateTokenMiddleware(handlers.Home)).Methods("GET")
    r.HandleFunc("/addBook", auth.ValidateTokenMiddleware(handlers.AddBook)).Methods("POST")
    r.HandleFunc("/deleteBook", auth.ValidateTokenMiddleware(handlers.DeleteBook)).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", r))
}
