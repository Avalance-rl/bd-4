package standard

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *sql.DB
}

type PGX struct {
    poll *pgxpool.Pool
}

func InitDB(opts string) *Storage {
	db, err := sql.Open("postgres", opts)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{db: db}
}

func InitDBWithPool(opts string, maxConns int32) *PGX {
    config, err := pgxpool.ParseConfig(opts)
	if err != nil {
        log.Fatal(err)
	}
	config.MaxConns = maxConns

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
        log.Fatal(err)
	}
	if err = db.Ping(context.Background()); err != nil {
        log.Fatal(err)
	}

	return &PGX{db} 
}

func (s *Storage) InsertProduct(p *Product) error {
	if p == nil {
		return nil
	}
	query := `
        INSERT INTO products (price, title, description)
        VALUES ($1, $2, $3)
        RETURNING id`

	err := s.db.QueryRow(query, p.Price, p.Title, p.Description).Scan(&p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SelectProducts() ([]Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Price, &p.Title, &p.Description); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *Storage) GetProduct(id string) (*Product, error) {
	if id == "" {
		return nil, nil
	}

	var product Product
	err := s.db.QueryRow("SELECT * FROM products WHERE id = $1", id).Scan(&product)
	if err != nil {
		return nil, nil
	}

	return &product, nil
}

func (s *Storage) UpdateProduct(p *Product) error {
	query := `
        UPDATE products 
        SET price = $1, title = $2, description = $3
        WHERE id = $4`

	result, err := s.db.Exec(query, p.Price, p.Title, p.Description, p.ID)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) InsertProducts(products []Product) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT INTO products (price, title, description)
        VALUES ($1, $2, $3)
        RETURNING id`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for i := range products {
		err = stmt.QueryRow(
			products[i].Price,
			products[i].Title,
			products[i].Description,
		).Scan(&products[i].ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
