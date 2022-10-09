package model

import (
	"encoding/json"
)

type User struct {
	UserId    string   `gorm:"primaryKey" json:"userId"`
	Email     string   `gorm:"column:email" json:"userEmail"`
	Name      string   `gorm:"column:name" json:"userName"`
	CreatedAt string   `gorm:"column:created_at"`
	IsActive  bool     `gorm:"type:boolean; column:is_active"`
	UserRole  UserRole `json:"userRole"`
	Pic       Pic      `json:"picDetail"`
}

func (User) TableName() string {
	return "u_user"
}

func (u *User) ToJSON() string {
	user, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		return ""
	}
	return string(user)
}
