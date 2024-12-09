package db

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// InitDB initializes a PostgreSQL database connection.
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// InsertProduct inserts a new product into the database.
func InsertProduct(db *sql.DB, userID int, name, description string, images []string, price float64) (int, error) {
	var productID int
	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(query, userID, name, description, images, price).Scan(&productID)
	return productID, err
}
