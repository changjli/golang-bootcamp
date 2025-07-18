package main

import (
	transaction "core-service/domains/transaction/entities"
	"core-service/domains/users/entities"
	wallet "core-service/domains/wallet/entities"
	"fmt"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	fmt.Println("Running Migrations...")

	// GORM will create tables, add missing columns, and update indexes.
	// It will not delete unused columns to protect your data.
	err := db.AutoMigrate(
		&entities.User{},
		&wallet.Wallet{},
		&transaction.Transaction{},
	)

	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database migrated successfully!")
}
