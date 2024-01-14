// handlers/handlers.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"uas-pemplat/models"

	"github.com/gorilla/mux"
)

// Handler struct contains handlers for User and Product operations
type Handler struct {
	Db *sql.DB
}

// NewHandler creates a new Handler with the given database connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{Db: db}
}

// CreateUser handles the creation of a new user
// @Summary Create a new user
// @Description Create a new user with the provided details
// @ID create-user
// @Accept json
// @Produce json
// @Param user body User true "User details"
// @Success 201 {object} map[string]int64 "user_id"
// @Failure 400 {object} ErrorResponse
// @Router /user [post]

// CreateUser handles the creation of a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := models.CreateUser(h.Db, &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"user_id": userID})
}

// GetUserByID retrieves a user by ID
// @Summary Get user by ID
// @Description Get user details by ID
// @ID get-user-by-id
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User "user"
// @Failure 400 {object} ErrorResponse
// @Router /user/{id} [get]

// GetUserByID retrieves a user by ID
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(h.Db, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// GetUsers retrieves all users
// @Summary Get all users
// @Description Get details of all users
// @ID get-all-users
// @Produce json
// @Success 200 {array} User "users"
// @Failure 500 {object} ErrorResponse
// @Router /users [get]

// GetUsers retrieves all users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers(h.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving users: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// UpdateUser updates a user by ID
// @Summary Update user by ID
// @Description Update user details by ID
// @ID update-user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body User true "Updated user details"
// @Success 200 {string} OK "user updated"
// @Failure 400 {object} ErrorResponse
// @Router /user/{id} [put]

// UpdateUser updates a user by ID
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = models.UpdateUser(h.Db, userID, user.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser deletes a user by ID
// @Summary Delete user by ID
// @Description Delete user by ID
// @ID delete-user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} OK "user deleted"
// @Failure 400 {object} ErrorResponse
// @Router /user/{id} [delete]

// DeleteUser deletes a user by ID
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteUser(h.Db, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and provide access token
// @ID user-login
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {string} OK "Login successful"
// @Failure 401 {object} ErrorResponse
// @Router /login [post]

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the LoginUser function to perform user authentication
	success, err := models.LoginUser(h.Db, loginRequest.Username, loginRequest.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Login failed: %v"}`, err), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if success {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"success": true, "message": "Login successful"}`)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, `{"success": false, "error": "Invalid username or password"}`)
	}
}

// CreateProduct handles the creation of a new product
// @Summary Create a new product
// @Description Create a new product with the provided details
// @ID create-product
// @Accept json
// @Produce json
// @Param product body Product true "Product details"
// @Success 201 {object} map[string]int64 "product_id"
// @Failure 400 {object} ErrorResponse
// @Router /product [post]

// CreateProduct handles the creation of a new product
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	productID, err := models.CreateProduct(h.Db, &product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"product_id": productID})
}

// GetProductByID retrieves a product by ID
// @Summary Get product by ID
// @Description Get product details by ID
// @ID get-product-by-id
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product "product"
// @Failure 400 {object} ErrorResponse
// @Router /product/{id} [get]

// GetProductByID retrieves a product by ID
func (h *Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(h.Db, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving product: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// GetProducts retrieves all products
// @Summary Get all products
// @Description Get details of all products
// @ID get-all-products
// @Produce json
// @Success 200 {array} Product "products"
// @Failure 500 {object} ErrorResponse
// @Router /products [get]

// GetProducts retrieves all products
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := models.GetProducts(h.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving products: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// UpdateProduct updates a product by ID
// @Summary Update product by ID
// @Description Update product details by ID
// @ID update-product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Updated product details"
// @Success 200 {string} OK "product updated"
// @Failure 400 {object} ErrorResponse
// @Router /product/{id} [put]

// UpdateProduct updates a product by ID
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = models.UpdateProduct(h.Db, productID, product.ProductName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct deletes a product by ID
// @Summary Delete product by ID
// @Description Delete product by ID
// @ID delete-product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {string} OK "product deleted"
// @Failure 400 {object} ErrorResponse
// @Router /product/{id} [delete]

// DeleteProduct deletes a product by ID
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteProduct(h.Db, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RenderLogin renders the login page
// @Summary Render login page
// @Description Render the login page HTML
// @ID render-login
// @Produce html
// @Success 200 {string} OK "login page rendered"
// @Failure 500 {object} ErrorResponse
// @Router /render-login [get]

// RenderLogin renders the login page
func RenderLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the login HTML template
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderHome renders the home page
// @Summary Render home page
// @Description Render the home page HTML
// @ID render-home
// @Produce html
// @Success 200 {string} OK "home page rendered"
// @Failure 500 {object} ErrorResponse
// @Router /render-home [get]

// RenderHome renders the home page
func RenderHome(w http.ResponseWriter, r *http.Request) {
	// Parse the home HTML template
	tmpl, err := template.ParseFiles("views/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
