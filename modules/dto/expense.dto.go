package dto

import "time"

type ExpenseDto struct {
	CreatedAt     int     `json:"createdAt"`
	Branch        string  `json:"branch" binding:"required"`
	ExpenseDetail string  `json:"expenseDetail" binding:"required"`
	ExpenseCost   float32 `json:"expenseCost" binding:"required"`
}

type ExpenseList struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Branch        string    `json:"branch" binding:"required"`
	ExpenseDetail string    `json:"expenseDetail" binding:"required"`
	ExpenseCost   float32   `json:"expenseCost" binding:"required"`
}

type ExpenseReport struct {
	SummaryCost  float32        `json:"summaryCost"`
	Transactions *[]ExpenseList `json:"transactions"`
}
