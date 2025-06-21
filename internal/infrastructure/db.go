package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the GORM DB connection.
type DB struct {
	conn *gorm.DB
}

// Config holds the database configuration.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDB creates a new database connection.
func NewDB(config Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	conn, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection.
	sqlDB, dbErr := conn.DB()
	if dbErr != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", dbErr)
	}
	pingErr := sqlDB.Ping()
	if pingErr != nil {
		return nil, fmt.Errorf("failed to ping database: %w", pingErr)
	}

	return &DB{
		conn: conn,
	}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	closeErr := sqlDB.Close()
	if closeErr != nil {
		return fmt.Errorf("failed to close database connection: %w", closeErr)
	}
	return nil
}

// Exec executes given SQL.
func (db *DB) Exec(sql string, values ...any) error {
	return db.conn.Exec(sql, values...).Error
}

// GetDB returns the underlying *gorm.DB instance.
func (db *DB) GetDB() *gorm.DB {
	return db.conn
}
