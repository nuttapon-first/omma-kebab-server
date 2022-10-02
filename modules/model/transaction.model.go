package model

import "time"

type Transaction struct {
	ID                 int       `gorm:"primarykey" json:"id"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time
	MenuID             int     `json:"menuId" binding:"required"`
	Branch             string  `json:"branch" binding:"required"`
	TransactionType    string  `json:"transactionType" binding:"required"`
	Channel            string  `json:"channel" binding:"required"`
	TransactionPrice   int     `json:"transactionPrice" binding:"required"`
	TransactionUnit    int     `json:"transactionUnit" binding:"required"`
	Fee                float32 `json:"fee"`
	Vat                float32 `json:"vat"`
	Discount           float32 `json:"discount"`
	TotalPrice         float32 `json:"totalPrice"`
	TotalCost          float32 `json:"totalCost"`
	TotalProfit        float32 `json:"totalProfit"`
	TotalProfitPercent float32 `json:"totalProfitPercent"`
	PaymentChannel     string  `json:"paymentChannel" binding:"required"`
	AddOn              string  `json:"addOn"`
}

func (Transaction) TableName() string {
	return "Transactions"
}
