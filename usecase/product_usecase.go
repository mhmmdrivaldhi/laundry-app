package usecase

import (
	"fmt"
	"go-laundry-app/models"
	repositories "go-laundry-app/repositories"
)

type productUsecase struct {
	repo repositories.ProductRepository
}

type ProductUsecase interface {
	CreateNewProduct(product models.Product) (models.Product, error)
	GetAllProduct() ([]models.Product,  error)
	GetProductByID(id int) (models.Product, error)
	UpdateProductByID(product models.Product) (models.Product, error)
	DeleteProductByID(id int) error 
}

func (p *productUsecase) CreateNewProduct(product models.Product) (models.Product, error) {
	return p.repo.CreateNewProduct(product)
}

func (p *productUsecase) GetAllProduct() ([]models.Product, error) {
	return p.repo.GetAllProduct()
}

func (p *productUsecase) GetProductByID(id int) (models.Product, error) {
	return p.repo.GetProductByID(id)
}

func (p *productUsecase) UpdateProductByID(product models.Product) (models.Product, error) {
	_,  err := p.repo.GetProductByID(product.ID)
	if err != nil {
		return models.Product{}, fmt.Errorf("product With ID %d Not Found", product.ID)
	}
	return p.repo.UpdateProductByID(product)
}

func (p *productUsecase) DeleteProductByID(id int) error {
	_, err := p.repo.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("product With ID %d Not Found", id)
	}
	return p.repo.DeleteProductByID(id)
}

func NewProductUsecase(repo repositories.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}