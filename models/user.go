package models

import "time"

type User struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Name      string `json:"name"`
	EmailID   string `json:"email_id"`
	Password  string `json:"password"`
	PhoneNo   string `json:"phone_no"`
}
