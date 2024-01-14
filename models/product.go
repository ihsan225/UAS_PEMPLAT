// models/product.go
package models

import (
	"database/sql"
	"time"
)

// Product struct represents the Products table in the database
type Product struct {
	ProductID     int       `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateProduct function inserts a new product into the Products table
func CreateProduct(db *sql.DB, product *Product) (int64, error) {
	result, err := db.Exec("INSERT INTO products (product_name, description, price, stock_quantity) VALUES (?, ?, ?, ?)",
		product.ProductName, product.Description, product.Price, product.StockQuantity)
	if err != nil {
		return 0, err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return productID, nil
}

// GetProductByID function retrieves a product from the Products table based on product_id
func GetProductByID(db *sql.DB, productID int) (*Product, error) {
	product := &Product{}
	err := db.QueryRow("SELECT product_id, product_name, description, price, stock_quantity, created_at FROM products WHERE product_id = ?", productID).
		Scan(&product.ProductID, &product.ProductName, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProducts function retrieves all products from the Products table
func GetProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT product_id, product_name, description, price, stock_quantity, created_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductID, &product.ProductName, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct function updates a product in the Products table
func UpdateProduct(db *sql.DB, productID int, newProductName string) error {
	_, err := db.Exec("UPDATE products SET product_name = ? WHERE product_id = ?", newProductName, productID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct function deletes a product from the Products table
func DeleteProduct(db *sql.DB, productID int) error {
	_, err := db.Exec("DELETE FROM products WHERE product_id = ?", productID)
	if err != nil {
		return err
	}

	return nil
}
