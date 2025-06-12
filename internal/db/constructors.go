// Package db contains the database schema and persistence layer
package db

// NewUserObject returns a new User with default values
func NewUserObject() User {
	// I will return only the fields that are not omitted
	return User{
		Language: "en",
		Units:    "metric",
		IsActive: true,
	}
}
