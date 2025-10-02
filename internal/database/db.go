package database

import (
	"fmt"
	"log"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	
	// Run auto-migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	// Check and create tables only if they don't exist
	tables := []struct {
		model interface{}
		name  string
	}{
		{&models.Company{}, "companies"},
		{&models.User{}, "users"},
		{&models.Project{}, "projects"},
		{&models.Lead{}, "leads"},
		{&models.Deal{}, "deals"},
		{&models.Proposal{}, "proposals"},
	}
	
	for _, table := range tables {
		if !db.Migrator().HasTable(table.name) {
			log.Printf("Creating table: %s", table.name)
			if err := db.AutoMigrate(table.model); err != nil {
				log.Printf("Error creating table %s: %v", table.name, err)
			}
		} else {
			log.Printf("Table already exists: %s (skipping)", table.name)
		}
	}
	
	log.Println("Database migrations completed")
	return nil
}
