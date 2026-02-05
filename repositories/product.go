package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAllProducts(name string) ([]*models.ProductWithCategory, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name AS category_name 
				FROM products AS p 
				JOIN categories AS c ON p.category_id = c.id`

	var args []interface{}
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*models.ProductWithCategory, 0)
	for rows.Next() {
		var p models.ProductWithCategory
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (repo *ProductRepository) GetAllProductsByCategoryID(categoryID int) ([]*models.ProductWithCategory, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name AS category_name 
				FROM products AS p JOIN categories AS c ON p.category_id = c.id 
				WHERE p.category_id = $1`
	rows, err := repo.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*models.ProductWithCategory, 0)
	for rows.Next() {
		var p models.ProductWithCategory
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (repo *ProductRepository) CreateProduct(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetProductByID(id int) (*models.ProductWithCategory, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name AS category_name 
				FROM products AS p JOIN categories AS c ON p.category_id = c.id 
				WHERE p.id = $1`

	var p models.ProductWithCategory
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("Produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) UpdateProduct(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Produk tidak ditemukan")
	}

	return err
}
