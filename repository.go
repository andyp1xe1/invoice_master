package main

import (
	"gorm.io/gorm"
	"log"
)

// CreateContract creates a Contract with associated Services
func CreateContract(db *gorm.DB, contract *Contract) error {
	err := db.Create(contract).Error
	if err != nil {
		log.Println("Error creating contract:", err)
		return err
	}
	return nil
}

// GetContractByID retrieves a Contract by ID with associated Services
func GetContractByID(db *gorm.DB, contractID uint) (*Contract, error) {
	var contract Contract
	err := db.Preload("Services").First(&contract, contractID).Error
	if err != nil {
		log.Println("Error fetching contract:", err)
		return nil, err
	}
	return &contract, nil
}
