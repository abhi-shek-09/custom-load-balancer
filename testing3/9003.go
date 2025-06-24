package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Simulate processing time
		time.Sleep(500 * time.Millisecond)
		fmt.Fprintf(w, "Hello from 9003!")
	})

	http.ListenAndServe(":9003", nil)
}
