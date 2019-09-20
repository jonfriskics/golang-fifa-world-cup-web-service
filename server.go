package main

import (
	"golang-fifa-world-cup-web-service/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/winners", handlers.WinnersHandler)
	http.ListenAndServe(":8000", nil)
}
