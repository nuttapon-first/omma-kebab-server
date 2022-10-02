package model

import "time"

type Stock struct {
	ID          int `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Ingredient  string  `json:"ingredient" binding:"required"`
	StockCost   float32 `json:"stockCost" binding:"required"`
	StockAmount float32 `json:"stockAmount" binding:"required"`
}

type AddStock struct {
	Ingredient  string  `json:"ingredient" binding:"required"`
	StockAmount float32 `json:"stockAmount" binding:"required"`
}

func (Stock) TableName() string {
	return "Stocks"
}
