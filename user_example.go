package main

import (
	"errors"
	"strings"
)

// User represents a user in the system
type User struct {
	ID       int
	Username string
	Email    string
	Active   bool
}

// UserStore manages a collection of users
type UserStore struct {
	users  map[int]User
	nextID int
}

// AddUser adds a new user and returns the created user or an error
func (us *UserStore) AddUser(username, email string) (User, error) {
	// Validate input
	if username == "" {
		return User{}, errors.New("username is required")
	}
	if !strings.Contains(email, "@") {
		return User{}, errors.New("email is invalid")
	}

	// Create new user
	user := User{
		ID:       us.nextID,
		Username: username,
		Email:    email,
		Active:   true,
	}

	us.users[user.ID] = user
	us.nextID++

	return user, nil
}
