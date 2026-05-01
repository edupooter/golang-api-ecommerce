package repo

import (
	"database/sql"
	"fmt"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	_ "modernc.org/sqlite"
)

// SQLiteOrderRepo implements OrderRepository using SQLite
type SQLiteOrderRepo struct {
	db *sql.DB
}

// SQLiteCustomerRepo implements CustomerRepository using SQLite
type SQLiteCustomerRepo struct {
	db *sql.DB
}

func NewSQLiteOrderRepo(dsn string) (*SQLiteOrderRepo, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	// ensure orders and order_items tables
	schema := `CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        customer_id INTEGER NOT NULL,
        total REAL NOT NULL
    );`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("create orders table: %w", err)
	}
	schema2 := `CREATE TABLE IF NOT EXISTS order_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        order_id INTEGER NOT NULL,
        product_id INTEGER NOT NULL,
        quantity INTEGER NOT NULL,
        price REAL NOT NULL
    );`
	if _, err := db.Exec(schema2); err != nil {
		db.Close()
		return nil, fmt.Errorf("create order_items table: %w", err)
	}
	return &SQLiteOrderRepo{db: db}, nil
}

func NewSQLiteCustomerRepo(dsn string) (*SQLiteCustomerRepo, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	schema := `CREATE TABLE IF NOT EXISTS customers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL
    );`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("create customers table: %w", err)
	}
	return &SQLiteCustomerRepo{db: db}, nil
}

func (r *SQLiteOrderRepo) Close() error {
	return r.db.Close()
}

func (r *SQLiteCustomerRepo) Close() error {
	return r.db.Close()
}

func (r *SQLiteOrderRepo) Create(o *model.Order) (*model.Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec("INSERT INTO orders(customer_id, total) VALUES(?, ?)", o.CustomerID, o.Total)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	oid, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, it := range o.Items {
		if _, err := tx.Exec("INSERT INTO order_items(order_id, product_id, quantity, price) VALUES(?, ?, ?, ?)", oid, it.ProductID, it.Quantity, it.Price); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	o.ID = oid
	cp := *o
	return &cp, nil
}

func (r *SQLiteCustomerRepo) GetByID(id int64) (*model.Customer, error) {
	var c model.Customer
	row := r.db.QueryRow("SELECT id, name, email FROM customers WHERE id = ?", id)
	if err := row.Scan(&c.ID, &c.Name, &c.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	cp := c
	return &cp, nil
}

func (r *SQLiteCustomerRepo) Create(c *model.Customer) (*model.Customer, error) {
	res, err := r.db.Exec("INSERT INTO customers(name, email) VALUES(?, ?)", c.Name, c.Email)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	c.ID = id
	cp := *c
	return &cp, nil
}
