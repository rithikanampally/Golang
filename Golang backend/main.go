package main

import (
	"backend/api"
	"backend/queue"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Import PostgreSQL driver
)

func main() {
	// Initialize DB connection
	dbConn, err := sql.Open("postgres", "postgres://postgres:Pass@123@localhost:5432/products_db?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	// Initialize RabbitMQ connection
	queueConn, channel, err := queue.InitRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer queueConn.Close()
	defer channel.Close()

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/products", api.CreateProductHandler(dbConn, channel)).Methods("POST")
	router.HandleFunc("/products/{id}", api.GetProductByIDHandler(dbConn)).Methods("GET")
	router.HandleFunc("/products", api.GetAllProductsHandler(dbConn)).Methods("GET")

	// Start server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
