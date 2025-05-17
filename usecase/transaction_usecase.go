package usecase

import (
	"fmt"
	"go-laundry-app/models"
	"go-laundry-app/repositories"
)

type transactionUsecase struct {
	repo repositories.TransactionRepository
}

type TransactionUsecase interface {
	CreateNewTransaction(transaction models.Transaction) (models.Transaction, error) 
	GetAllTransaction() ([]models.Transaction, error)
	GetTransactionByID(id int) (models.Transaction, error)
	UpdateTransactionByID(transaction models.Transaction) (models.Transaction, error)
	DeleteTransactionByID(id int) error 
}

func (t *transactionUsecase) CreateNewTransaction(transaction models.Transaction) (models.Transaction, error) {
	return t.repo.CreateNewTransaction(transaction)
}

func (t *transactionUsecase) GetAllTransaction() ([]models.Transaction, error) {
	return t.repo.GetAllTransaction()
}

func (t *transactionUsecase) GetTransactionByID(id int) (models.Transaction, error) {
	return t.repo.GetTransactionByID(id)
} 

func (t *transactionUsecase) UpdateTransactionByID(transaction models.Transaction) (models.Transaction, error) {
	_, err := t.repo.GetTransactionByID(transaction.ID)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("transaction with ID %d Not Found", transaction.ID)
	}
	return t.repo.UpdateTransactionByID(transaction)
}

func (t *transactionUsecase) DeleteTransactionByID(id int) error {
	_, err := t.repo.GetTransactionByID(id)
	if err != nil {
		return fmt.Errorf("transaction with ID %d Not Found", id)
	}
	return t.repo.DeleteTransactionByID(id)
}

func NewTransactionUsecase(repo repositories.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{repo: repo}
}
