package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/thejixer/shop-api/internal/models"
)

func (s *PostgresStore) createAddressTable() error {
	query := `create table if not exists address (
		id SERIAL PRIMARY KEY,
		userId integer,
		title VARCHAR(50),
		lon decimal,
		CHECK (lon>=-180 AND lon <= 180),
		lat decimal,
		CHECK (lat>=-90 AND lat <=90),
		address text,
		recieverName VARCHAR(50),
		recieverPhone VARCHAR(20),
		deleted bool DEFAULT false
	)`

	_, err := s.db.Exec(query)

	return err
}

type AddressRepo struct {
	db *sql.DB
}

func NewAddressRepo(db *sql.DB) *AddressRepo {
	return &AddressRepo{
		db: db,
	}
}

func (r *AddressRepo) Create(a *models.Address) (*models.Address, error) {

	query := `
		INSERT INTO address (userId, title, lon, lat, address, recieverName, recieverPhone)
		values ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	lastInsertId := 0

	insertErr := r.db.QueryRow(query, a.UserId, a.Title, a.Lon, a.Lat, a.Address, a.RecieverName, a.RecieverPhone).Scan(&lastInsertId)
	if insertErr != nil {
		fmt.Println(insertErr)
		return nil, insertErr
	}

	if lastInsertId == 0 {
		return nil, errors.New("some sort of problem with database")
	}

	a.Id = lastInsertId

	return a, nil
}

func (r *AddressRepo) FindById(id int) (*models.Address, error) {
	rows, err := r.db.Query("SELECT * FROM address WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoAddress(rows)
	}

	return nil, errors.New("not found")
}

func (r *AddressRepo) FindByUserId(userId int) ([]*models.Address, error) {
	query := "SELECT * FROM address WHERE userId = $1 ORDER BY id desc"
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	addresses := []*models.Address{}
	for rows.Next() {
		u, err := ScanIntoAddress(rows)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, u)
	}

	return addresses, nil
}

func (r *AddressRepo) Delete(id int) error {
	query := `
		UPDATE address
		SET deleted = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, true, id)
	if err != nil {
		return err
	}
	return nil

}

func ScanIntoAddress(rows *sql.Rows) (*models.Address, error) {
	a := new(models.Address)
	if err := rows.Scan(
		&a.Id,
		&a.UserId,
		&a.Title,
		&a.Lon,
		&a.Lat,
		&a.Address,
		&a.RecieverName,
		&a.RecieverPhone,
		&a.Deleted,
	); err != nil {
		return nil, err
	}
	return a, nil
}
