package auth

import (
	"time"
)

type AuthRequest struct {
	UserName string `json:"user_name"`
}

type AuthResponse struct {
	UserName string `json:"user_name"`
	Token    string `json:"token"`
	UserCode string `json:"user_code"`
}

type User struct {
	ID        int64     `gorm:"index;column:id;primaryKey" json:"id"`
	UserCode  string    `gorm:"index;type:uuid;column:user_code;not null" json:"user_code"`
	UserName  string    `gorm:"column:user_name;not null" json:"user_name"`
	TokenHash string    `gorm:"column:token_hash;not null" json:"token_hash"`
	Active    bool      `gorm:"column:active;not null" json:"active"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
