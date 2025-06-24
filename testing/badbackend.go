package main

// import (
// 	"log"
// 	"net/http"
// )

// func main() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("❌ Backend 9001: sending 500 error")
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 	})
// 	log.Println("🔥 Bad backend running on :9001")
// 	http.ListenAndServe(":9001", nil)
// }
