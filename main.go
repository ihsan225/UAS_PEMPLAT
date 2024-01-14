// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"uas-pemplat/handlers"

	// "uas-pemplat/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	// "github.com/yourusername/yourprojectname/handlers"
	// "github.com/yourusername/yourprojectname/models"
)

func main() {
	// Replace with your actual database connection details
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/uas-pemplat-2024?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database")

	// // Initialize the database connection in the models package
	// models.InitDB()
	// defer models.CloseDB()

	// Create instances of UserHandler and ProductHandler
	userHandler := handlers.NewHandler(db)
	productHandler := handlers.NewHandler(db)

	// Create a new router
	router := mux.NewRouter()

	// Define Login endpoint
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Define User endpoints
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Define Product endpoints
	router.HandleFunc("/products", productHandler.GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", productHandler.GetProductByID).Methods("GET")
	router.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Serve login page
	router.HandleFunc("/login", handlers.RenderLogin).Methods("GET")

	// Serve home page
	router.HandleFunc("/home", handlers.RenderHome).Methods("GET")

	// Start the HTTP server
	port := ":8083"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
