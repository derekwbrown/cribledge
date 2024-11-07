package main


import (
	"fmt"
	"net/http"
)


func main() {
	
	// set up a default http handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received a / request:", r.URL.Path)
		fmt.Println("GET params were:", r.URL.Query())

		// respond with a simple 200
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/bar/bar/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received a bar bar request:", r.URL.Path)
		fmt.Println("GET params were:", r.URL.Query())

		// respond with a simple 200
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/bar/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received a bar request:", r.URL.Path)
		fmt.Println("GET params were:", r.URL.Query())

		// respond with a simple 200
		w.WriteHeader(http.StatusOK)
	})
	
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received a foo request:", r.URL.Path)
		fmt.Println("GET params were:", r.URL.Query())

		// respond with a simple 200
		w.WriteHeader(http.StatusOK)
	})
	// start the server
	http.ListenAndServe(":8080", nil)
}