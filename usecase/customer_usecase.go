package usecase

import (
	"fmt"
	"go-laundry-app/models"
	"go-laundry-app/repositories"
)

type customerUsecase struct {
	repo repositories.CustomerRepository
}

type CustomerUsecase interface {
	CreateNewCustomer(customer models.Customer) (models.Customer, error)
	GetAllCustomer() ([]models.Customer, error)
	GetCustomerById(id int) (models.Customer, error)
	UpdateCustomerById(customer models.Customer) (models.Customer, error)
	DeleteCustomerById(id int) error
}

func (c *customerUsecase) CreateNewCustomer(customer models.Customer) (models.Customer, error) {
	return c.repo.CreateNewCustomer(customer)
}

func (c *customerUsecase) GetAllCustomer() ([]models.Customer, error) {
	return c.repo.GetAllCustomer()
}

func (c *customerUsecase) GetCustomerById(id int) (models.Customer, error) {
	return c.repo.GetCustomerById(id)
}

func (c *customerUsecase) UpdateCustomerById(customer models.Customer) (models.Customer, error) {
	_, err := c.repo.GetCustomerById(customer.ID)
	if err != nil {
		return models.Customer{}, fmt.Errorf("customer with ID %d Not Found", customer.ID)
	}
	return c.repo.UpdateCustomerById(customer)
}

func (c *customerUsecase) DeleteCustomerById(id int) error {
	_, err := c.repo.GetCustomerById(id)
	if err != nil {
		return fmt.Errorf("customer with ID %d Not Found", id)
	}
	return c.repo.DeleteCustomerById(id)
}

func NewCustomerUseCase(repo repositories.CustomerRepository) CustomerUsecase {
	return &customerUsecase{repo: repo}
}
