package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/thejixer/shop-api/internal/models"
)

func (s *PostgresStore) createCartItemTable() error {
	query := `create table if not exists cartitems (
		id SERIAL PRIMARY KEY,
		userId integer,
		productId integer,
		quantity integer
	)`

	_, err := s.db.Exec(query)

	return err
}

type CartRepo struct {
	db *sql.DB
}

func NewCartRepo(db *sql.DB) *CartRepo {
	return &CartRepo{
		db: db,
	}
}

func (r *CartRepo) Add(userId, productId, quantity int) error {

	// there should be a way to execute both these at the same db trip
	r.db.Exec("DELETE FROM cartitems WHERE cartitems.userId = $1 AND cartitems.productId = $2;", userId, productId)
	query := `
		INSERT INTO cartitems (userId, productId, quantity)
		VALUES ($1, $2, $3) RETURNING id`
	lastInsertId := 0

	insertErr := r.db.QueryRow(query, userId, productId, quantity).Scan(&lastInsertId)
	if insertErr != nil {
		return insertErr
	}

	if lastInsertId == 0 {
		return errors.New("some sort of problem with database")
	}

	return nil
}

func (r *CartRepo) FindUsersItems(userId int) ([]*models.CartItem, error) {

	query := `
		Select c.quantity as quantity, p.id as ProductId, p.title as ProductTitle, p.Price as ProductPrice, p.quantity as ProductQuantity, p.description as ProductDescription
		From cartitems c
		JOIN products p 
		ON c.productId = p.id
		WHERE c.userId = $1
	`
	rows, err := r.db.Query(query, userId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	items := []*models.CartItem{}
	for rows.Next() {
		u, err := scanIntoCartItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	return items, nil

}

func (r *CartRepo) Remove(userId, productId int) error {
	x, err := r.db.Exec("DELETE FROM cartitems WHERE cartitems.userId = $1 AND cartitems.productId = $2;", userId, productId)

	fmt.Println(x)
	if err != nil {
		return err
	}
	return nil
}

func scanIntoCartItem(rows *sql.Rows) (*models.CartItem, error) {
	x, _ := rows.Columns()
	fmt.Println("#############")
	fmt.Printf("%+v \n", x)
	fmt.Println("#############")
	c := new(models.CartItem)
	if err := rows.Scan(
		&c.Quantity,
		&c.ProductId,
		&c.ProductTitle,
		&c.ProductPrice,
		&c.ProductQuantity,
		&c.ProductDescription,
	); err != nil {
		return nil, err
	}
	return c, nil
}
