package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thejixer/shop-api/internal/models"
)

func (s *PostgresStore) createProductTable() error {
	query := `create table if not exists products (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100),
		price DECIMAL,
		quantity integer,
		description VARCHAR(1000),
		createdAt TIMESTAMP
	)`

	_, err := s.db.Exec(query)

	return err
}

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) Create(data models.CreateProductDto) (*models.Product, error) {

	n := &models.Product{
		Title:       data.Title,
		Price:       decimal.NewFromFloat(data.Price),
		Quantity:    data.Quantity,
		Description: data.Description,
		CreatedAt:   time.Now().UTC(),
	}
	query := `
	INSERT INTO PRODUCTS (title, price, quantity, description, createdAt)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	lastInsertId := 0
	insertErr := r.db.QueryRow(query, n.Title, n.Price, n.Quantity, n.Description, n.CreatedAt).Scan(&lastInsertId)
	if insertErr != nil {
		return nil, insertErr
	}
	n.Id = lastInsertId

	return n, nil
}

func (r *ProductRepo) Edit(id int, data models.CreateProductDto) (*models.Product, error) {
	query := `
		UPDATE PRODUCTS
		SET title = $2,
			price = $3,
			quantity = $4,
			description = $5
			WHERE id = $1
		RETURNING *;
	`

	rows, err := r.db.Query(query, id, data.Title, data.Price, data.Quantity, data.Description)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoProducts(rows)
	}

	return nil, errors.New("server error")

}

func (r *ProductRepo) FindById(id int) (*models.Product, error) {
	rows, err := r.db.Query("SELECT * FROM PRODUCTS WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoProducts(rows)
	}

	return nil, errors.New("not found")
}

func (r *ProductRepo) Find(text string, page, limit int) ([]*models.Product, int, error) {

	offset := page * limit
	query := `SELECT * FROM PRODUCTS 
		WHERE LOWER(PRODUCTS.title) LIKE $2 
		OR LOWER(PRODUCTS.description) LIKE $2 
		ORDER BY id
		OFFSET $1 ROWS
		FETCH NEXT $3 ROWS ONLY`
	str := "%" + strings.ToLower(text) + "%"
	rows, err := r.db.Query(query, offset, str, limit)
	// defer rows.Close()
	if err != nil {
		return nil, 0, err
	}
	products := []*models.Product{}
	for rows.Next() {
		u, err := scanIntoProducts(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, u)
	}

	var count int
	r.db.QueryRow(`
		SELECT count(id) 
		FROM PRODUCTS
		WHERE LOWER(PRODUCTS.title) LIKE $1
		OR LOWER(PRODUCTS.description) LIKE $1
	`, str).Scan(&count)

	return products, count, nil
}

func scanIntoProducts(rows *sql.Rows) (*models.Product, error) {
	x, _ := rows.Columns()
	fmt.Println("#############")
	fmt.Printf("%+v \n", x)
	fmt.Println("#############")
	p := new(models.Product)
	if err := rows.Scan(&p.Id, &p.Title, &p.Price, &p.Quantity, &p.Description, &p.CreatedAt); err != nil {
		return nil, err
	}
	return p, nil
}
