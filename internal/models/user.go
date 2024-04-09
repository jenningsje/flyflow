package models

type User struct {
	BaseModel
	Email          string `json:"email" gorm:"uniqueIndex"`
	HashedPassword string `json:"-"`
	Password       string `json:"-" gorm:"-"`
}