package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/thejixer/shop-api/models"
)

type PostgresStore struct {
	db       *sql.DB
	UserRepo models.UserRepository
}

func NewPostgresStore() (*PostgresStore, error) {

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	conString := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable", dbUser, dbName, dbPassword)
	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	userRepo := NewUserRepo(db)

	return &PostgresStore{
		db:       db,
		UserRepo: userRepo,
	}, nil
}

func (s *PostgresStore) Init() error {

	if err := s.createUserTable(); err != nil {
		return err
	}

	return nil

}
