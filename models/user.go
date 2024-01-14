// models/user.go
package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents the Users table in the database
type User struct {
	UserID           int       `json:"user_id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"password_hash"`
	RegistrationDate time.Time `json:"registration_date"`
}

// CreateUser function inserts a new user into the Users table
func CreateUser(db *sql.DB, user *User) (int64, error) {
	result, err := db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByID function retrieves a user from the Users table based on user_id
func GetUserByID(db *sql.DB, userID int) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT user_id, username, email, password_hash, registration_date FROM users WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.Username, &user.Email, &user.PasswordHash, &user.RegistrationDate)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers function retrieves all users from the Users table
func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT user_id, username, email, password_hash, registration_date FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.PasswordHash, &user.RegistrationDate)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser function updates a user in the Users table
func UpdateUser(db *sql.DB, userID int, newUsername string) error {
	_, err := db.Exec("UPDATE users SET username = ? WHERE user_id = ?", newUsername, userID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser function deletes a user from the Users table
func DeleteUser(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM users WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	return nil
}

// LoginUser authenticates a user and returns true if successful
func LoginUser(db *sql.DB, username, password string) (bool, error) {
	// Retrieve hashed password from the database
	var hashedPassword string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		// Passwords match, authentication successful
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		// Passwords don't match, authentication failed
		return false, nil
	} else {
		// Other error (e.g., invalid hash)
		return false, err
	}
}
