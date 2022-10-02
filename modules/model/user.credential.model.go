package model

type UserCredential struct {
	ID         int    `gorm:"primarykey" json:"id"`
	UserId     int    `json:"userId" binding:"required"`
	Credential string `json:"credential" binding:"required"`
}

func (UserCredential) TableName() string {
	return "UserCredentials"
}
