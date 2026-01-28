package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAllCategories() ([]*models.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (repo *CategoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var p models.Category
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *CategoryRepository) CreateCategory(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}

func (repo *CategoryRepository) UpdateCategory(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	_, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	return err
}

func (repo *CategoryRepository) DeleteCategory(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := repo.db.Exec(query, id)
	return err
}
