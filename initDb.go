package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

// initDb initializes the SQLite database and migrates the models
func initDb() (*gorm.DB, error) {
	// Open SQLite connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the models to create tables
	err = db.AutoMigrate(&Invoice{}, &Service{}, &Contract{})
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized and tables migrated.")
	return db, nil
}

