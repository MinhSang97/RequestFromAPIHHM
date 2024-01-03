package dbutil

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"

	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	url := oracle.BuildUrl("118.69.35.119", 1521, "hhm", "MiniMDM10", "MiniMDM10", nil)
	db, err := gorm.Open(oracle.Open(url), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	fmt.Println("connected to the database:", db)
	return db, err
}
