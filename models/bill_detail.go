package models

type BillDetail struct {
	ID           int  `json:"id"`
	BillID       int  `json:"billId"`
	ProductID    int  `json:"productId"`
	ProductPrice int    `json:"productPrice"`
	Qty          int     `json:"qty"`
	Product      *Product `json:"product,omitempty"`
}