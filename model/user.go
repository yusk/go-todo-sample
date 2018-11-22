package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primary_key;unique;index;not null;"`
	Email     string    `json:"email" gorm:"unique;index;not null;"`
	Password  string    `json:"-" gorm:"not null;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
