package dto

import (
	"time"
)

type TransactionDto struct {
	CreatedAt        int      `json:"createdAt"`
	MenuID           int      `json:"menuId" binding:"required"`
	Branch           string   `json:"branch" binding:"required"`
	TransactionType  string   `json:"transactionType" binding:"required"`
	Channel          string   `json:"channel" binding:"required"`
	TransactionPrice int      `json:"transactionPrice" binding:"required"`
	TransactionUnit  int      `json:"transactionUnit" binding:"required"`
	Fee              float32  `json:"fee"`
	Vat              float32  `json:"vat"`
	Discount         float32  `json:"discount"`
	PaymentChannel   string   `json:"paymentChannel" binding:"required"`
	AddOns           []AddOns `gorm:"-" json:"addOns"`
}

type AddOns struct {
	StockId int     `json:"stockId"`
	Amount  float32 `json:"amount"`
}

type TransactionList struct {
	ID                 int       `json:"id"`
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
	MenuName           string  `json:"menuName"`
}

type SalesReport struct {
	SummaryAmount    int                `json:"summaryAmount"`
	SummaryFee       float32            `json:"summaryFee"`
	SummaryVat       float32            `json:"summaryVat"`
	SummarySales     float32            `json:"summarySales"`
	SummaryCost      float32            `json:"summaryCost"`
	SummaryProfit    float32            `json:"summaryProfit"`
	AvgProfitPercent float32            `json:"avgProfitPercent"`
	Transactions     *[]TransactionList `json:"transactions"`
}
