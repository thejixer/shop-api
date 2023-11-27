package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/thejixer/shop-api/internal/models"
)

func (s *PostgresStore) createOrderItemTable() error {
	query := `
	create table if not exists orders (
		id SERIAL PRIMARY KEY,
		userId integer,
		status VARCHAR(16),
		totalPrice DECIMAL,
		createdAt TIMESTAMP
	);

	create table if not exists orderItems (
		id SERIAL PRIMARY KEY,
		orderId integer,
		productId int,
		quantity int,
		price DECIMAL
	)	
	
	`
	_, err := s.db.Exec(query)

	return err
}

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func getUsersCartItems(db *sql.DB, userId int) ([]*models.CartItem, error) {
	query := `
		Select c.quantity as quantity, p.id as ProductId, p.title as ProductTitle, p.Price as ProductPrice, p.quantity as ProductQuantity, p.description as ProductDescription
		From cartitems c
		JOIN products p 
		ON c.productId = p.id
		WHERE c.userId = $1
	`
	rows, err := db.Query(query, userId)
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

func (r *OrderRepo) Create(data models.Order) error {

	userCartItems, err := getUsersCartItems(r.db, data.UserId)
	if err != nil {
		return err
	}
	if len(userCartItems) == 0 {
		return errors.New("there is nothing in your cart")
	}
	var totalPrice decimal.Decimal

	for _, s := range userCartItems {
		Quantity := decimal.NewFromInt(int64(s.Quantity))
		tp := Quantity.Mul(s.ProductPrice)
		totalPrice = totalPrice.Add(tp)
	}

	query := fmt.Sprintf(`
		BEGIN TRANSACTION;

		UPDATE USERS
		SET balance = balance - %v
		WHERE id = %v;

		INSERT INTO orders (userId, status, totalPrice, createdAt)
		values(%v, '%v', %v, Now());

	`, totalPrice, data.UserId, data.UserId, data.Status, totalPrice)

	for _, s := range userCartItems {
		query += fmt.Sprintf(`INSERT INTO orderItems (orderId, productId, quantity, price)
		values (currval('orders_id_seq'), %v, %v, %v);
		
		`, s.ProductId, s.Quantity, s.ProductPrice)
		query += fmt.Sprintf(`DELETE FROM cartitems WHERE cartitems.userId = %v AND cartitems.productId = %v;
		
		`, data.UserId, s.ProductId)
		query += fmt.Sprintf(`
			UPDATE products
			SET quantity = quantity - %v
			WHERE id = %v;

		`, s.Quantity, s.ProductId)
	}

	query += `
		COMMIT;
	`

	_, transactionError := r.db.Exec(query)

	if transactionError != nil {
		return transactionError
	}

	return nil
}

func (r *OrderRepo) FindById(id int) (*models.Order, error) {
	rows, err := r.db.Query("SELECT * FROM orders WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoOrders(rows)
	}

	return nil, errors.New("not found")
}

func (r *OrderRepo) FindOrderItemsOfSingleOrder(orderId, userId int) ([]*models.OrderItem, error) {
	query := `
		Select i.orderId as orderId, i.quantity as quantity, i.price as priceAtTheTime, p.id as ProductId, p.title as ProductTitle, p.Price as ProductPrice, p.quantity as ProductQuantity, p.description as ProductDescription
		From orderItems i
		JOIN products p
		ON i.productId = p.id
		WHERE i.orderId = $1
		`
	rows, err := r.db.Query(query, orderId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	items := []*models.OrderItem{}
	for rows.Next() {
		u, err := scanIntoItems(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, u)
	}

	return items, nil
}

func scanIntoOrders(rows *sql.Rows) (*models.Order, error) {
	o := new(models.Order)
	if err := rows.Scan(&o.Id, &o.UserId, &o.Status, &o.TotalPrice, &o.CreatedAt); err != nil {
		return nil, err
	}
	return o, nil
}

func scanIntoItems(rows *sql.Rows) (*models.OrderItem, error) {
	i := new(models.OrderItem)
	if err := rows.Scan(
		&i.OrderId,
		&i.Quantity,
		&i.PriceAtTheTime,
		&i.ProductId,
		&i.ProductTitle,
		&i.ProductPrice,
		&i.ProductQuantity,
		&i.ProductDescription,
	); err != nil {
		return nil, err
	}
	return i, nil
}
