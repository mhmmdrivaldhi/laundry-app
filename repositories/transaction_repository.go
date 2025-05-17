package repositories

import (
	"database/sql"
	"go-laundry-app/models"
	"log"
)

type transactionRepository struct {
	db *sql.DB
}
 
type TransactionRepository interface {
	CreateNewTransaction(transaction models.Transaction) (models.Transaction, error) 
	GetAllTransaction() ([]models.Transaction, error)
	GetTransactionByID(id int) (models.Transaction, error)
	UpdateTransactionByID(transaction models.Transaction) (models.Transaction, error)
	DeleteTransactionByID(id int) error 
	GetEmployeeByID(id int) (models.Employee, error)
	GetCustomerByID(id int) (models.Customer, error)
	GetProductByID(id int) (models.Product, error)
	AddBillDetail(detail models.BillDetail) (models.BillDetail, error)
	DeleteBillDetailByTransactionID(transaction_id int) error

}

func (t *transactionRepository) CreateNewTransaction(transaction models.Transaction) (models.Transaction, error) {
	var transaction_id int
	total_bill := 0

	for i, detail := range transaction.BillDetails {
		product, err := t.GetProductByID(detail.ProductID)
		if err != nil {
			log.Println("error getting product", err)
			return models.Transaction{}, err
		}
		transaction.BillDetails[i].ProductPrice = product.Price
		total_bill += product.Price * detail.Qty
	}

	err := t.db.QueryRow("INSERT INTO transactions(bill_date, entry_date, finish_date, employee_id, customer_id, total_bill) VALUES($1, $2, $3, $4, $5, $6) RETURNING id", transaction.BillDate, transaction.EntryDate, transaction.FinishDate, transaction.EmployeeID, transaction.CustomerID, total_bill).Scan(&transaction_id)

	if err != nil {
		return models.Transaction{}, err
	}

	for i, bill_detail := range transaction.BillDetails {
		bill_detail.BillID = transaction_id
		new_detail, err := t.AddBillDetail(bill_detail)
		if err != nil {
			log.Println("error insert bill details", err)
			return models.Transaction{}, err
		}
		transaction.BillDetails[i] = new_detail
	}
	transaction.ID = transaction_id
	transaction.TotalBill = total_bill
	return transaction, nil
}

func (t *transactionRepository) GetAllTransaction() ([]models.Transaction, error) {
	var transactions []models.Transaction

	rows, err := t.db.Query("SELECT id, bill_date, entry_date, finish_date, employee_id, customer_id, total_bill FROM transactions")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(&transaction.ID, &transaction.BillDate, &transaction.EntryDate, &transaction.FinishDate, &transaction.EmployeeID, &transaction.CustomerID, &transaction.TotalBill)

		if err != nil {
			return nil, err
		}

		employee, err := t.GetEmployeeByID(transaction.EmployeeID)
		if err != nil {
			log.Println("error getting employee:", err)
			return nil, err
		}
		transaction.Employee = employee

		customer, err := t.GetCustomerByID(transaction.CustomerID)
		if err != nil {
			log.Println("error getting customer:", err)
			return nil, err
		}
		transaction.Customer = customer

		bill_details, err := t.db.Query("SELECT id, bill_id, product_id, product_price, qty FROM bill_details WHERE bill_id = $1", transaction.ID)
		if err != nil {
			log.Println("error getting bill details:", err)
			return nil, err
		}

		defer bill_details.Close()

		for bill_details.Next() {
			 var bill_detail models.BillDetail
			 err := bill_details.Scan(&bill_detail.ID, &bill_detail.BillID, &bill_detail.ProductID, &bill_detail.ProductPrice, &bill_detail.Qty)
			 if err != nil {
				log.Println("error Scanning bill details:", err)
				return nil, err
			 }

			 product, err := t.GetProductByID(bill_detail.ProductID)
			 if err != nil {
				log.Println("error getting Product:", err)
				return nil, err
			 }
			 bill_detail.Product = &product

			 transaction.BillDetails = append(transaction.BillDetails, bill_detail)
		}

		transactions = append(transactions, transaction)
	}
	return transactions, nil
} 

func (t *transactionRepository) GetTransactionByID(id int) (models.Transaction, error) {
	var transaction models.Transaction

	tx_query := `SELECT id, bill_date, entry_date, finish_date, employee_id, customer_id, total_bill FROM transactions WHERE id = $1`

	err := t.db.QueryRow(tx_query, id).Scan(
		&transaction.ID, &transaction.BillDate, &transaction.EntryDate,
		&transaction.FinishDate, &transaction.EmployeeID, &transaction.CustomerID, &transaction.TotalBill)
	if err != nil {
		log.Println("error:", err)
		return models.Transaction{}, err
	}

	employee, err := t.GetEmployeeByID(transaction.EmployeeID)
	if err != nil {
		log.Println("error:", err)
		return models.Transaction{}, err
	}
	transaction.Employee = employee
	
	customer, err := t.GetCustomerByID(transaction.CustomerID)
	if err != nil {
		log.Println("error:", err)
		return models.Transaction{}, err
	}
	transaction.Customer = customer

	rows, err := t.db.Query(`SELECT id, bill_id, product_id, product_price, qty FROM bill_details WHERE bill_id = $1`, transaction.ID)
	if err != nil {
		log.Println("error:", err)
		return models.Transaction{}, err
	}
	
	defer rows.Close()

	for rows.Next() {
		bill_detail := models.BillDetail{}

		err := rows.Scan(&bill_detail.ID, &bill_detail.BillID, &bill_detail.ProductID, &bill_detail.ProductPrice, &bill_detail.Qty)
		if err != nil {
			log.Println("error:", err)
			return models.Transaction{}, err
		}

		// Get Nested Product 	
		product, err := t.GetProductByID(bill_detail.ProductID)
		if err != nil {
			log.Println("error:", err)
			return models.Transaction{}, err
		}
		bill_detail.Product = &product

		transaction.BillDetails = append(transaction.BillDetails, bill_detail)
	}
	return transaction, nil
}

func (t *transactionRepository) UpdateTransactionByID(transaction models.Transaction) (models.Transaction, error) {
	err := t.DeleteBillDetailByTransactionID(transaction.ID)
	if err != nil {
		log.Println("error deleting old bill details:", err)
		return models.Transaction{}, err
	}

	total_bill := 0
	for i, bill_detail := range transaction.BillDetails {
		product, err := t.GetProductByID(bill_detail.ProductID)
		if err != nil {
			log.Println("error getting product")
			return models.Transaction{}, err
		}
		transaction.BillDetails[i].ProductPrice = product.Price
		transaction.BillDetails[i].BillID = transaction.ID
		total_bill += product.Price * bill_detail.Qty

		new_detail, err := t.AddBillDetail(transaction.BillDetails[i])
		if err != nil {
			log.Println("error inserting bill detail:", err)
			return models.Transaction{}, err
		}
		transaction.BillDetails[i] = new_detail
	}

	_, err = t.db.Exec("UPDATE transactions SET bill_date = $1, entry_date = $2, finish_date = $3, employee_id = $4, customer_id = $5, total_bill = $6 WHERE id = $7",transaction.BillDate, transaction.EntryDate, transaction.FinishDate, transaction.EmployeeID, transaction.CustomerID, transaction.TotalBill, transaction.ID )
	if err != nil {
		return models.Transaction{}, err
	}
	transaction.TotalBill = total_bill
	return transaction, nil
}

func (t *transactionRepository) DeleteTransactionByID(id int) error {
	_, err := t.db.Exec("DELETE FROM transactions WHERE id = $1",id)
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionRepository) GetEmployeeByID(id int) (models.Employee, error) {
	employee := models.Employee{}

	err := t.db.QueryRow("SELECT id, name, phone_number, address FROM employees WHERE id = $1",id).Scan(&employee.ID, &employee.Name, &employee.PhoneNumber, &employee.Address)
	if err != nil {
		log.Println("error:", err)
	}
	return employee, nil
}

func (t *transactionRepository) GetCustomerByID(id int) (models.Customer, error) {
	customer := models.Customer{}

	err := t.db.QueryRow("SELECT id, name, phone_number, address FROM customers WHERE id = $1",id).Scan(&customer.ID, &customer.Name, &customer.PhoneNumber, &customer.Address)
	if err != nil {
		log.Println("error:", err)
	}
	return customer, nil
}

func (t *transactionRepository) GetProductByID(id int) (models.Product, error) {
	product := models.Product{}

	err := t.db.QueryRow("SELECT id, name, price, unit FROM products WHERE id = $1",id).Scan(&product.ID, &product.Name, &product.Price, &product.Unit)
	if err != nil {
		log.Println("error:", err)
	}
	return product, nil
}

func (t *transactionRepository) AddBillDetail(detail models.BillDetail) (models.BillDetail, error) {
	_, err := t.db.Exec("INSERT INTO bill_details(bill_id, product_id, product_price, qty) VALUES($1, $2, $3, $4)", detail.BillID, detail.ProductID, detail.ProductPrice, detail.Qty)
	if err != nil {
		return models.BillDetail{}, err
	}
	return detail, nil
} 

func (t *transactionRepository) DeleteBillDetailByTransactionID(id int) error {
	_, err := t.db.Exec("DELETE FROM bill_details WHERE bill_id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}