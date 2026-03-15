package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disables pgx driver prepared statements
	}), &gorm.Config{
		PrepareStmt: false, // Disables GORM prepared statements
	})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}
