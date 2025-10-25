package initializers

import (
	"errors"
	"task_2/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect to the MySQL database
func ConnectToDB(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connectionString))
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Perform migrations
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return errors.New("Database connection can't be nil")
	}
	err := db.AutoMigrate(&models.Country{})
	if err != nil {
		return err
	}
	return nil
}