package model

import "time"

type Expense struct {
	ID            int       `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time
	Branch        string  `json:"branch" binding:"required"`
	ExpenseDetail string  `json:"expenseDetail" binding:"required"`
	ExpenseCost   float32 `json:"expenseCost" binding:"required"`
}

func (Expense) TableName() string {
	return "Expenses"
}
