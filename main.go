package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

func main() {
	fmt.Println("Starting server...")

	http.HandleFunc("/health-check", healthCheckHandler)

	s := &http.Server{}
	log.Fatal(s.ListenAndServe())
}
