package repositories

import (
	"database/sql"
	"go-laundry-app/models"
)

type customerRepository struct {
	db *sql.DB
}

type CustomerRepository interface {
	CreateNewCustomer(customer models.Customer) (models.Customer, error)
	GetAllCustomer() ([]models.Customer, error)
	GetCustomerById(id int) (models.Customer, error)
	UpdateCustomerById(customer models.Customer) (models.Customer, error)
	DeleteCustomerById(id int) error
}

func (c *customerRepository) CreateNewCustomer(customer models.Customer) (models.Customer, error) {
	var customer_id int
	err := c.db.QueryRow("INSERT INTO customers(name, phone_number, address) VALUES($1, $2, $3) RETURNING id", customer.Name, customer.PhoneNumber, customer.Address).Scan(&customer_id)

	if err != nil {
		return models.Customer{}, err
	}
	customer.ID = customer_id
	return customer, nil
}

func (c *customerRepository) GetAllCustomer() ([]models.Customer, error) {
	var customers []models.Customer

	rows, err := c.db.Query("SELECT id, name, phone_number, address FROM customers")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var customer models.Customer

		err := rows.Scan(&customer.ID, &customer.Name, &customer.PhoneNumber, &customer.Address)

		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (c *customerRepository) GetCustomerById(id int) (models.Customer, error) {
	var customer models.Customer

	err := c.db.QueryRow("SELECT id, name, phone_number, address FROM customers WHERE id = $1", id).Scan(&customer.ID, &customer.Name, &customer.PhoneNumber, &customer.Address)

	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (c *customerRepository) UpdateCustomerById(customer models.Customer) (models.Customer, error) {
	_, err := c.db.Exec("UPDATE customers SET name = $2, phone_number = $3, address = $4  WHERE id = $1", customer.ID, customer.Name, customer.PhoneNumber, customer.Address)

	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (c *customerRepository) DeleteCustomerById(id int) error {
	_, err := c.db.Exec("DELETE FROM customers WHERE id = $1", id)

	if err != nil {
		return err
	}
	return nil
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}
