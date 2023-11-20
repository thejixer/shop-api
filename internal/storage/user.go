package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/thejixer/shop-api/internal/models"
	"github.com/thejixer/shop-api/pkg/encryption"
)

func (s *PostgresStore) createUserTable() error {

	query := `create table if not exists users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100),
		isEmailVerified BOOLEAN,
		password VARCHAR,
		balance DECIMAL,
		createdAt TIMESTAMP
	)`

	_, err := s.db.Exec(query)

	return err
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(data models.SignUpDTO) (*models.User, error) {

	hashedPassword, err := encryption.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Name:      data.Name,
		Email:     data.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	query := `
	INSERT INTO USERS (name, email, isEmailVerified, password, balance, createdAt)
	VALUES ($1, LOWER($2), $3, $4, $5, $6) RETURNING id`
	lastInsertId := 0

	insertErr := r.db.QueryRow(query, newUser.Name, newUser.Email, false, newUser.Password, 0, newUser.CreatedAt).Scan(&lastInsertId)
	if insertErr != nil {
		return nil, insertErr
	}
	newUser.ID = lastInsertId

	return newUser, nil
}
func (r *UserRepo) FindById(id int) (*models.User, error) {
	rows, err := r.db.Query("SELECT * FROM USERS WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoUsers(rows)
	}

	return nil, errors.New("not found")
}
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	rows, err := r.db.Query("SELECT * FROM USERS WHERE email = LOWER($1)", email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoUsers(rows)
	}

	return nil, errors.New("not found")
}
func (r *UserRepo) VerifyEmail(email string) error {
	query := `
		UPDATE USERS
		SET isEmailVerified = $2
		WHERE email = LOWER($1);
	`
	_, err := r.db.Exec(query, email, true)

	if err != nil {
		return err
	}

	return nil
}
func (r *UserRepo) UpdatePassword(email, password string) error {
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return err
	}

	query := `
		UPDATE USERS
		SET password = $2
		WHERE email = LOWER($1);
	`
	_, updateErr := r.db.Exec(query, email, hashedPassword)

	if updateErr != nil {
		return err
	}

	return nil
}

func scanIntoUsers(rows *sql.Rows) (*models.User, error) {
	u := new(models.User)
	if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsEmailVerified, &u.Password, &u.Balance, &u.CreatedAt); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return u, nil
}
