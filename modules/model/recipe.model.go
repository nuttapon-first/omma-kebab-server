package model

type Recipe struct {
	ID               int     `gorm:"primarykey" json:"id"`
	MenuID           int     `json:"menuId"`
	StockID          int     `json:"stockId" binding:"required"`
	IngredientAmount float32 `json:"ingredientAmount" binding:"required"`
}

func (Recipe) TableName() string {
	return "Recipes"
}
