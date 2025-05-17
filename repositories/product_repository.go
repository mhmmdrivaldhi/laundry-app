package repositories

import (
	"database/sql"
	"go-laundry-app/models"
)

type productRepository struct {
	db *sql.DB
}

type ProductRepository interface {
	CreateNewProduct(product models.Product) (models.Product, error)
	GetAllProduct() ([]models.Product, error)
	GetProductByID(id int) (models.Product, error)
	UpdateProductByID(product models.Product) (models.Product, error)
	DeleteProductByID(id int) error 
}

func (p *productRepository) CreateNewProduct(product models.Product) (models.Product, error) {
	var product_id int

	err := p.db.QueryRow("INSERT INTO products(name, price, unit) VALUES($1, $2, $3) RETURNING id", product.Name, product.Price, product.Unit).Scan(&product_id)

	if err != nil {
		return models.Product{}, err
	}
	product.ID = product_id
	return product, nil
}

func (p *productRepository) GetAllProduct() ([]models.Product, error) {
	var products []models.Product

	rows, err := p.db.Query("SELECT id, name, price, unit FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Unit)

		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *productRepository) GetProductByID(id int) (models.Product, error) {
	var product models.Product

	err := p.db.QueryRow("SELECT id, name, price, unit FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price, &product.Unit)

	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (p *productRepository) UpdateProductByID(product models.Product) (models.Product, error) {
	_, err := p.db.Exec("UPDATE products SET name = $2, price = $3, unit = $4 WHERE  id = $1", product.ID, product.Name, product.Price, product.Unit)

	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (p *productRepository) DeleteProductByID(id int) error {
	_, err := p.db.Exec("DELETE FROM products WHERE id = $1", id)

	if err != nil {
		return err
	}
	return nil
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}