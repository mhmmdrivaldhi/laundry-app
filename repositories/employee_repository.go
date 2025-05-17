package repositories

import (
	"database/sql"
	"go-laundry-app/models"
)

type employeeRepository struct {
	db *sql.DB
}

type EmployeeRepository interface {
	CreateNewEmployee(employee models.Employee) (models.Employee, error)
	GetAllEmployee() ([]models.Employee, error)
	GetEmployeeById(id int) (models.Employee, error)
	UpdateEmployeeById(employee models.Employee) (models.Employee, error)
	DeleteEmployeeById(id int) error
}

func (e *employeeRepository) CreateNewEmployee(employee models.Employee) (models.Employee, error) {
	var employee_id int
	err := e.db.QueryRow("INSERT INTO employees(name, phone_number, address) VALUES($1, $2, $3) RETURNING id", employee.Name, employee.PhoneNumber, employee.Address).Scan(&employee_id)

	if err != nil {
		return models.Employee{}, err
	}
	employee.ID = employee_id
	return employee, nil
}

func (e *employeeRepository) GetAllEmployee() ([]models.Employee, error) {
	var employees []models.Employee

	rows, err := e.db.Query("SELECT id, name, phone_number, address FROM employees")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var employee models.Employee

		err := rows.Scan(&employee.ID, &employee.Name, &employee.PhoneNumber, &employee.Address)

		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (e *employeeRepository) GetEmployeeById(id int) (models.Employee, error) {
	var employee models.Employee

	err := e.db.QueryRow("SELECT id, name, phone_number, address FROM employees WHERE id = $1", id).Scan(&employee.ID, &employee.Name, &employee.PhoneNumber, &employee.Address)

	if err != nil {
		return models.Employee{}, err
	}
	return employee, nil
}

func (e *employeeRepository) UpdateEmployeeById(employee models.Employee) (models.Employee, error) {
	_, err := e.db.Exec("UPDATE employees SET name = $2, phone_number = $3, address = $4  WHERE id = $1", employee.ID, employee.Name, employee.PhoneNumber, employee.Address)

	if err != nil {
		return models.Employee{}, err
	}
	return employee, nil
}

func (e *employeeRepository) DeleteEmployeeById(id int) error {
	_, err := e.db.Exec("DELETE FROM employees WHERE id = $1", id)

	if err != nil {
		return err
	}
	return nil
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}