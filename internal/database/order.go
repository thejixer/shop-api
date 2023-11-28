package database

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/shopspring/decimal"
	dataprocesslayer "github.com/thejixer/shop-api/internal/data-process-layer"
	"github.com/thejixer/shop-api/internal/models"
)

func (s *PostgresStore) createOrderItemTable() error {
	query := `
	create table if not exists orders (
		id SERIAL PRIMARY KEY,
		userId integer,
		addressId integer,
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

func (r *OrderRepo) Create(data models.Order, cartItems []*models.CartItem) error {

	var totalPrice decimal.Decimal

	for _, s := range cartItems {
		x := (decimal.NewFromInt(int64(s.Quantity))).Mul(s.ProductPrice)
		totalPrice = totalPrice.Add(x)
	}

	txn, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, insertErr := txn.Exec(`
		UPDATE USERS
		SET balance = balance - $1
		WHERE id = $2;
	`, totalPrice, data.UserId)

	if insertErr != nil {
		return insertErr
	}

	_, insertOrderErr := txn.Exec(`
		INSERT INTO orders (userId, addressId, status, totalPrice, createdAt)
		values($1, $2, $3, $4, Now());
	`, data.UserId, data.AddressId, data.Status, totalPrice)

	if insertOrderErr != nil {
		return insertOrderErr
	}

	for _, s := range cartItems {

		_, err := txn.Exec(`INSERT INTO orderItems (orderId, productId, quantity, price)
		values (currval('orders_id_seq'), $1, $2, $3);`, s.ProductId, s.Quantity, s.ProductPrice)
		if err != nil {
			return err
		}

		_, err2 := txn.Exec(`	
			UPDATE products
			SET quantity = quantity - $1
			WHERE id = $2;`, s.Quantity, s.ProductId)

		if err2 != nil {
			return err2
		}

	}

	_, cartDelErr := txn.Exec(`DELETE FROM cartitems WHERE cartitems.userId = $1`, data.UserId)
	if cartDelErr != nil {
		return cartDelErr
	}

	commitErr := txn.Commit()
	if commitErr != nil {
		return commitErr
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

func (r *OrderRepo) MakeOrder(order *models.Order, t *models.OrderDto, e *error, wg *sync.WaitGroup) {

	defer wg.Done()

	var err error
	var rows *sql.Rows
	// get User
	var thisUser *models.User
	rows, err = r.db.Query("SELECT * FROM USERS WHERE id = $1", order.UserId)
	if err != nil {
		*e = err
		return
	}
	for rows.Next() {
		user, err := ScanIntoUsers(rows)
		if err != nil {
			*e = err
			return
		}
		thisUser = user
	}
	// get Address
	user := dataprocesslayer.ConvertToUserDto(thisUser)
	var thisAddress *models.Address
	rows, err = r.db.Query("SELECT * FROM address WHERE id = $1", order.AddressId)
	if err != nil {
		*e = err
		return
	}
	for rows.Next() {
		address, err := ScanIntoAddress(rows)
		if err != nil {
			*e = err
			return
		}
		thisAddress = address
	}
	address := dataprocesslayer.ConvertToAddressDto(thisAddress, user)
	// get order items
	rows, err = r.db.Query(`
	Select i.orderId as orderId, i.quantity as quantity, i.price as priceAtTheTime, p.id as ProductId, p.title as ProductTitle, p.Price as ProductPrice, p.quantity as ProductQuantity, p.description as ProductDescription
	From orderItems i
	JOIN products p
	ON i.productId = p.id
	WHERE i.orderId = $1
	`, order.Id)
	if err != nil {
		*e = err
		return
	}
	rawItems := []*models.OrderItem{}
	for rows.Next() {
		u, err := scanIntoItems(rows)
		if err != nil {
			*e = err
			return
		}
		rawItems = append(rawItems, u)
	}

	var items []models.OrderItemDto

	var temp decimal.Decimal

	for _, e := range rawItems {
		items = append(items, *dataprocesslayer.ConvertOrderItemtoDto(e))
		Quantity := decimal.NewFromInt(int64(e.Quantity))
		tp := Quantity.Mul(e.PriceAtTheTime)
		temp = temp.Add(tp)
	}
	totalPrice, _ := temp.Float64()

	t.Id = order.Id
	t.User = user
	t.Address = *address
	t.Status = order.Status
	t.Items = items
	t.CreatedAt = order.CreatedAt
	t.TotalPrice = totalPrice

}

func (r *OrderRepo) GetOrdersByUserId(userId, page, limit int) ([]*models.Order, int, error) {

	offset := page * limit
	query := "SELECT * FROM orders WHERE userId = $1 ORDER BY id desc OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY"
	rows, err := r.db.Query(query, userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	orders := []*models.Order{}
	for rows.Next() {
		u, err := scanIntoOrders(rows)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, u)
	}

	var count int
	r.db.QueryRow("SELECT count(id) FROM orders WHERE userId = $1", userId).Scan(&count)

	return orders, count, nil
}

func scanIntoOrders(rows *sql.Rows) (*models.Order, error) {
	o := new(models.Order)

	if err := rows.Scan(&o.Id, &o.UserId, &o.AddressId, &o.Status, &o.TotalPrice, &o.CreatedAt); err != nil {
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
