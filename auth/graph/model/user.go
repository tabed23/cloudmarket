package model

import "time"

type UserModel struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `gorm:"unique" json:"email"`
	Password     string    `json:"password"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type NewUserModel struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}
func (UserModel) TableName() string {
	return "users"
}