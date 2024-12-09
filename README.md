Golang Backend Project
This repository contains a backend system implemented in Golang. It demonstrates modular architecture and follows best practices for scalability, maintainability, and performance.

Features
RESTful API for managing resources (e.g., products, users).
Database integration using PostgreSQL.
Middleware for authentication, error handling, and request validation.
Structured logging with libraries like logrus.
Unit and integration testing for code reliability.
Setup Instructions
Prerequisites
Go (v1.16+ recommended)
PostgreSQL
Git
Steps
Clone the repository:

bash
Copy code
git clone https://github.com/rithikanampally/Golang.git
cd Golang
Install dependencies:

bash
Copy code
go mod tidy
Configure environment variables:

Create a .env file in the root directory.
Add configurations:
env
Copy code
DB_HOST=localhost
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
DB_PORT=5432
Run database migrations (if applicable):

bash
Copy code
go run migrate.go
Start the application:

bash
Copy code
go run main.go
Test the API using Postman or curl.

API Endpoints
Product Management
Create a Product
POST /products
Request body:

json
Copy code
{
    "user_id": 1,
    "product_name": "Sample Product",
    "product_description": "Description here",
    "product_price": 99.99
}
Get Product by ID
GET /products/:id

List Products
GET /products

Testing
Run unit tests with:

bash
Copy code
go test ./...
