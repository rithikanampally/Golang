package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"backend/queue"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/streadway/amqp"
)

func InsertProduct(db *sql.DB, userID int, name, description string, images []string, price float64) (int, error) {
	var productID int
	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// Use pq.Array to convert []string to a PostgreSQL array
	err := db.QueryRow(query, userID, name, description, pq.Array(images), price).Scan(&productID)
	if err != nil {
		log.Printf("Error inserting product: %v", err) // Log the error
		return 0, err
	}
	return productID, nil
}

type Product struct {
	ID                      int      `json:"id"`
	UserID                  int      `json:"user_id"`
	ProductName             string   `json:"product_name"`
	ProductDescription      string   `json:"product_description"`
	ProductImages           []string `json:"product_images"`
	CompressedProductImages []string `json:"compressed_product_images"`
	ProductPrice            float64  `json:"product_price"`
	CreatedAt               string   `json:"created_at"`
}

func GetAllProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		priceMin := r.URL.Query().Get("price_min")
		priceMax := r.URL.Query().Get("price_max")
		productName := r.URL.Query().Get("product_name")

		// Ensure that user_id is provided and is a valid integer
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		// Prepare the base query and arguments
		query := "SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price, created_at FROM products WHERE user_id = $1"
		args := []interface{}{userID}

		// Add price filtering if provided
		if priceMin != "" && priceMax != "" {
			query += " AND product_price BETWEEN $2 AND $3"
			args = append(args, priceMin, priceMax)
		}

		// Add product name filtering if provided
		if productName != "" {
			query += " AND LOWER(product_name) LIKE $4"
			args = append(args, "%"+productName+"%")
		}

		// Execute the query
		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving products: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Collect products from the query result
		var products []Product
		for rows.Next() {
			var product Product
			err := rows.Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription,
				pq.Array(&product.ProductImages), pq.Array(&product.CompressedProductImages), &product.ProductPrice, &product.CreatedAt)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading product data: %v", err), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}

		// Send the response as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}
func GetProductByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID := vars["id"]

		query := `SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price, created_at 
                  FROM products WHERE id = $1`

		var product Product
		err := db.QueryRow(query, productID).Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription,
			pq.Array(&product.ProductImages), pq.Array(&product.CompressedProductImages), &product.ProductPrice, &product.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Product not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error retrieving product: %v", err), http.StatusInternalServerError)
			}
			return
		}

		// Send the response as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

type CreateProductRequest struct {
	UserID      int      `json:"user_id"`
	Name        string   `json:"product_name"`
	Description string   `json:"product_description"`
	Images      []string `json:"product_images"`
	Price       float64  `json:"product_price"`
}

func CreateProductHandler(db *sql.DB, mqChannel *amqp.Channel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		productID, err := InsertProduct(db, req.UserID, req.Name, req.Description, req.Images, req.Price)
		if err != nil {
			http.Error(w, "Failed to insert product", http.StatusInternalServerError)
			return
		}

		// Publish image URLs to the queue
		imageData, err := json.Marshal(req.Images)
		if err != nil {
			http.Error(w, "Failed to process images", http.StatusInternalServerError)
			return
		}

		if err := queue.PublishToQueue(mqChannel, "image_processing", imageData); err != nil {
			http.Error(w, "Failed to enqueue images", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Product created successfully",
			"productID": productID,
		})
	}
}
