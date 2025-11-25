// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"go-auth/config"
	"go-auth/routes"
)

func main() {
	// Connect to database
	config.ConnectDB()

	// Register routes
	routes.UserRoutes()

	// Home Route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(" Server is running on port 8080"))
	})

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
