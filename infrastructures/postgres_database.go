package infrastructures

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *PostgresDatabase
)

func NewPostgresDatabase() Database {
	err := godotenv.Load()
	if err != nil {
		panic(".env file not found")
	}

	once.Do(func() {
		db, err := gorm.Open(postgres.Open(
			fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s",
				os.Getenv("HOST"),
				os.Getenv("USER"),
				os.Getenv("PASSWORD"),
				os.Getenv("DBNAME"),
				os.Getenv("PORT"),
				os.Getenv("SSLMODE"),
				os.Getenv("TIMEZONE"),
			),
		),
			&gorm.Config{},
		)

		if err != nil {
			panic(err)
		}

		dbInstance = &PostgresDatabase{
			db: db,
		}
	})

	return dbInstance
}

func (p *PostgresDatabase) GetInstance() *gorm.DB {
	return dbInstance.db
}
