package repo

import (
    "database/sql"
    "fmt"

    "github.com/edupooter/golang-api-ecommerce/internal/model"
    _ "modernc.org/sqlite"
)

type SQLiteRepo struct {
    db *sql.DB
}

func NewSQLiteRepo(dsn string) (*SQLiteRepo, error) {
    db, err := sql.Open("sqlite", dsn)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }
    // create table if not exists
    schema := `CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price REAL NOT NULL,
        stock INTEGER NOT NULL
    );`
    if _, err := db.Exec(schema); err != nil {
        db.Close()
        return nil, fmt.Errorf("create schema: %w", err)
    }
    // seed initial data if table empty
    var cnt int
    if err := db.QueryRow("SELECT COUNT(1) FROM products").Scan(&cnt); err == nil {
        if cnt == 0 {
            seed := []struct{
                name string
                price float64
                stock int
            }{
                {"Camiseta Golang", 49.9, 10},
                {"Caneca Go", 19.9, 25},
            }
            for _, s := range seed {
                _, _ = db.Exec("INSERT INTO products(name, price, stock) VALUES(?, ?, ?)", s.name, s.price, s.stock)
            }
        }
    }
    return &SQLiteRepo{db: db}, nil
}

func (r *SQLiteRepo) Close() error {
    return r.db.Close()
}

func (r *SQLiteRepo) Create(p *model.Product) (*model.Product, error) {
    res, err := r.db.Exec("INSERT INTO products(name, price, stock) VALUES(?, ?, ?)", p.Name, p.Price, p.Stock)
    if err != nil {
        return nil, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    cp := *p
    cp.ID = id
    return &cp, nil
}

func (r *SQLiteRepo) GetAll() ([]*model.Product, error) {
    rows, err := r.db.Query("SELECT id, name, price, stock FROM products")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var res []*model.Product
    for rows.Next() {
        var p model.Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
            return nil, err
        }
        cp := p
        res = append(res, &cp)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return res, nil
}

func (r *SQLiteRepo) GetByID(id int64) (*model.Product, error) {
    var p model.Product
    row := r.db.QueryRow("SELECT id, name, price, stock FROM products WHERE id = ?", id)
    if err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, err
    }
    cp := p
    return &cp, nil
}

func (r *SQLiteRepo) Update(p *model.Product) (*model.Product, error) {
    res, err := r.db.Exec("UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?", p.Name, p.Price, p.Stock, p.ID)
    if err != nil {
        return nil, err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return nil, err
    }
    if n == 0 {
        return nil, ErrNotFound
    }
    cp := *p
    return &cp, nil
}

func (r *SQLiteRepo) Delete(id int64) error {
    res, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
    if err != nil {
        return err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if n == 0 {
        return ErrNotFound
    }
    return nil
}
