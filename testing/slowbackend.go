package main

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("⚠️  Backend received request — simulating delay")
// 	time.Sleep(10 * time.Second) // Simulate slow processing
// 	w.Write([]byte("✅ Response from slow backend"))
// }

// func main() {
// 	http.HandleFunc("/", handler)
// 	fmt.Println("🚀 Slow backend running on http://localhost:9001")
// 	http.ListenAndServe(":9001", nil)
// }