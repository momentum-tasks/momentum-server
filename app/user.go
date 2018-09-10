package app

import (
	"golang.org/x/crypto/bcrypt"
)

// User is the structure to hold the information for a user, and all tasks and reports associated with them
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string

	Tasks []Task
}

var userTableStatement = `CREATE TABLE IF NOT EXISTS users (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	username VARCHAR(255) UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	passwordhash TEXT,
	PRIMARY KEY (id)
)`

// CreateUser creates a user with all the basic information required, and stores them in the database
func CreateUser(username string, email string, password string) error {
	hashed, err := HashPassword(password)
	if err != nil {
		return err
	}

	stmt, err := store.db.Prepare("INSERT INTO users(username, email, passwordhash) VALUES(?, ?, ?)")
	_, err = stmt.Exec(username, email, hashed)
	if err != nil {
		return err
	}
	return nil
}

// HashPassword generates a bcrypt hashed password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash verifies if the supplied plain-text password should match the stored hashed password for a user
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// GetUserByID queries the database for a user, and returns that user if they are found
// returns nil, with an error if the user is not found
func GetUserByID(ID int) (*User, error) {
	var u User
	err := store.db.QueryRow("SELECT id, username, email, passwordhash FROM users WHERE id = ?", ID).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	u.Tasks = GetTasksByUser(&u)
	return &u, nil
}

// GetUserByUsername queries the database for a user, and returns that user if they are found
// returns nil, with an error if the user is not found
func GetUserByUsername(username string) (*User, error) {
	var u User
	err := store.db.QueryRow("SELECT id, username, email, passwordhash FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
