package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// NewDatabase sets up and returns a new database connection
func NewDBConn(dsn string) (*Database, error) {
	// Set up the logger to log queries (optional)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// Configure connection pool (optional)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB from GORM: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate schema (you might want to do this elsewhere)
	if err := migrateModels(db); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	return &Database{DB: db}, nil
}
// MigrateModels migrates the models to the database
func migrateModels(d *gorm.DB) error {
	if err := d.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	if err := d.AutoMigrate(&model.Product{}); err != nil {
		return err
	}
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB from GORM: %v", err)
	}
	return sqlDB.Close()
}