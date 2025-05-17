package usecase

import (
	"fmt"
	"go-laundry-app/models"
	"go-laundry-app/repositories"
)

type employeeUsecase struct {
	repo repositories.EmployeeRepository
}

type EmployeeUsecase interface {
	CreateNewEmployee(employee models.Employee) (models.Employee, error)
	GetAllEmployee() ([]models.Employee, error)
	GetEmployeeById(id int) (models.Employee, error)
	UpdateEmployeeById(employee models.Employee) (models.Employee, error)
	DeleteEmployeeById(id int) error
}

func (e *employeeUsecase) CreateNewEmployee(employee models.Employee) (models.Employee, error) {
	return e.repo.CreateNewEmployee(employee)
}

func (e *employeeUsecase) GetAllEmployee() ([]models.Employee, error) {
	return e.repo.GetAllEmployee()
}

func (e *employeeUsecase) GetEmployeeById(id int) (models.Employee, error) {
	return e.repo.GetEmployeeById(id)
}

func (e *employeeUsecase) UpdateEmployeeById(employee models.Employee) (models.Employee, error) {
	_,  err := e.repo.GetEmployeeById(employee.ID)
	if err != nil {
		return models.Employee{}, fmt.Errorf("employee with ID %d Not Found", employee.ID)
	}
	return e.repo.UpdateEmployeeById(employee)
}

func (e *employeeUsecase) DeleteEmployeeById(id int) error {
	_, err := e.repo.GetEmployeeById(id)
	if err != nil {
		return fmt.Errorf("employee with ID %d Not Found", id)
	}
	return e.repo.DeleteEmployeeById(id)
}

func NewEmployeeUseCase(repo repositories.EmployeeRepository) EmployeeUsecase {
	return &employeeUsecase{repo: repo}
}

