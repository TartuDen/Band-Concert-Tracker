package main

import (
	"groupie-tracker/visualization"
	"net/http"
)



func main() {
	http.HandleFunc("/", visualization.HandlerMainPage)
	// http.HandleFunc("/", visualization.HandlerFilter)
	http.HandleFunc("/error", visualization.ErrorHandler)
	http.HandleFunc("/map", visualization.MapHandler)

	// Start the HTTP server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
