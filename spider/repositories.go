// Generated by OpenYan LazyAPI

package main

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)


// Cat represents the Cat model
type Cat struct {
	id uuid.UUID `json:"id"`
	breed string `json:"breed"`
	name string `json:"name"`
	
}

// InsertCat inserts a new Cat into the database
func InsertCat(db *sql.DB, id uuid.UUID, breed string, name string, ) error {
	query := "INSERT INTO Cat (id, breed, name) VALUES (?, ?, ?)"
	_, err := db.Exec(query, id, breed, name, )
	if err != nil {
		return fmt.Errorf("failed to insert Cat: %w", err)
	}
	return nil
}

