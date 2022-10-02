package model

type User struct {
	ID             int            `gorm:"primarykey" json:"id"`
	UserName       string         `gorm:"uniqueIndex" json:"userName" binding:"required"`
	UserFullName   string         `json:"userFullName" binding:"required"`
	UserRole       string         `json:"userRole" binding:"required"`
	UserCredential UserCredential `gorm:"ForeignKey:UserId;" json:"userCredential"`
}

type LoginUser struct {
	UserName     string `json:"username" binding:"required"`
	UserPassword string `json:"password" binding:"required"`
}

func (User) TableName() string {
	return "Users"
}
