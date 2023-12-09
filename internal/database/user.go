package database

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/thejixer/shop-api/internal/models"
	"github.com/thejixer/shop-api/internal/utils"
	"github.com/thejixer/shop-api/pkg/encryption"
)

func (s *PostgresStore) createUserTable() error {

	query := `
	create table if not exists users (
		id SERIAL PRIMARY KEY,
		role VARCHAR(50),
		name VARCHAR(100),
		email VARCHAR(100) UNIQUE,
		isEmailVerified BOOLEAN,
		password VARCHAR,
		balance DECIMAL,
		CHECK (balance>=0),
		permissions valid_permissions[],
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

func (r *UserRepo) Create(name, email, password, role string, isEmailVerified bool, permissions []string) (*models.User, error) {

	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Name:            name,
		Email:           email,
		Role:            role,
		Password:        hashedPassword,
		IsEmailVerified: isEmailVerified,
		Permissions:     permissions,
		CreatedAt:       time.Now().UTC(),
	}

	query := `
	INSERT INTO USERS (role, name, email, isEmailVerified, password, balance, permissions, createdAt)
	VALUES ($1, $2, LOWER($3), $4, $5, $6, $7, $8) RETURNING id`
	lastInsertId := 0

	insertErr := r.db.QueryRow(
		query,
		newUser.Role,
		newUser.Name,
		newUser.Email,
		newUser.IsEmailVerified,
		newUser.Password,
		0,
		pq.Array(newUser.Permissions),
		newUser.CreatedAt,
	).Scan(&lastInsertId)

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
		return ScanIntoUsers(rows)
	}

	return nil, errors.New("not found")
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	rows, err := r.db.Query("SELECT * FROM USERS WHERE email = LOWER($1)", email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoUsers(rows)
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
		return updateErr
	}

	return nil
}
func (r *UserRepo) FindUsers(text string, page, limit int) ([]*models.User, int, error) {

	offset := page * limit
	query := "SELECT * FROM USERS WHERE LOWER(USERS.name) LIKE $2 ORDER BY id OFFSET $1 ROWS FETCH NEXT $3 ROWS ONLY"
	str := "%" + strings.ToLower(text) + "%"
	rows, err := r.db.Query(query, offset, str, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	users := []*models.User{}
	for rows.Next() {
		u, err := ScanIntoUsers(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	var count int
	r.db.QueryRow("SELECT count(id) FROM USERS WHERE LOWER(USERS.name) LIKE $1", str).Scan(&count)

	return users, count, nil
}
func (r *UserRepo) ChargeBalance(userId int, amount float64) error {

	query := `
		UPDATE USERS
		SET balance = balance + $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, amount, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) UpdatePermissions(id int, permissions []string) error {

	permissions = utils.RemoveDuplicate[string](permissions)

	query := `
		UPDATE USERS
		SET permissions = $1
		WHERE (id = $2)
	`

	_, err := r.db.Exec(query, pq.Array(permissions), id)
	if err != nil {
		return err
	}
	return nil
}

func ScanIntoUsers(rows *sql.Rows) (*models.User, error) {
	u := new(models.User)
	if err := rows.Scan(
		&u.ID,
		&u.Role,
		&u.Name,
		&u.Email,
		&u.IsEmailVerified,
		&u.Password,
		&u.Balance,
		pq.Array(&u.Permissions),
		&u.CreatedAt,
	); err != nil {
		return nil, err
	}
	return u, nil
}
