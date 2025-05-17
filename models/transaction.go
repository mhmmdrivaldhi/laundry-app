package models

type Transaction struct {
	ID         int          `json:"id"`
	BillDate   string       `json:"billDate"`
	EntryDate  string       `json:"entryDate"`
	FinishDate string       `json:"finishDate"`

	
	EmployeeID int          `json:"employeeId"`
	CustomerID int          `json:"customerId"`
	
	Employee Employee 		 `json:"employee"`
	Customer Customer 		 `json:"customer"`
	BillDetails []BillDetail `json:"billDetails"`

	TotalBill   int          `json:"totalBill"`
}
