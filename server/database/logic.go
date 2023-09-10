package database

import (
	"DadJokesAPI/shared"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connection *gorm.DB

func connect() (*gorm.DB, error) {
	if connection == nil {
		user, err := shared.GetStr("postgres.user")
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres user: %w", err)
		}
		password, err := shared.GetStr("postgres.password")
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres password: %w", err)
		}
		dbName, err := shared.GetStr("postgres.db")
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres db name: %w", err)
		}
		port, err := shared.GetInt("postgres.port")
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres port: %w", err)
		}

		dsn := fmt.Sprintf("host=0.0.0.0 user=%s password=%s dbname=%s port=%d", user, password, dbName, port)
		dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		db, err := dbConnection.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get database connection: %w", err)
		}

		db.SetMaxIdleConns(shared.GetIntDef("postgres.max_idle_connections", 8))
		db.SetMaxOpenConns(shared.GetIntDef("postgres.max_open_connections", 64))

		if shared.GetBoolDef("postgres.query_debugging", false) {
			dbConnection = dbConnection.Debug()
		}
		connection = dbConnection
	}

	return connection, nil
}
