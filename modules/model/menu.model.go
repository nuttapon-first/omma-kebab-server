package model

import "time"

type Menu struct {
	ID               int `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	MenuName         string        `json:"menuName" binding:"required"`
	MenuNameTH       string        `json:"menuNameTH"`
	MenuUnit         int           `json:"menuUnit"`
	MenuCost         float32       `json:"menuCost" binding:"required"`
	MenuType         string        `json:"menuType" binding:"required"`
	MenuRecipe       []Recipe      `gorm:"ForeignKey:MenuID;" json:"menuRecipe"`
	MenuTransactions []Transaction `gorm:"ForeignKey:MenuID;"`
}

type EditMenu struct {
	ID               int `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	MenuName         string   `json:"menuName" binding:"required"`
	MenuNameTH       string   `json:"menuNameTH"`
	MenuUnit         int      `json:"menuUnit"`
	MenuCost         float32  `json:"menuCost" binding:"required"`
	MenuType         string   `json:"menuType" binding:"required"`
	MenuRecipe       []Recipe `gorm:"ForeignKey:MenuID;" json:"menuRecipe"`
	DeletedRecipeIds []string `gorm:"-" json:"deletedRecipeIds"`
}

func (Menu) TableName() string {
	return "MenuList"
}
