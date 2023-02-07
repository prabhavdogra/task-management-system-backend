package models

import "time"

type Task struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Heading   string `json:"heading"`
	Content   string `json:"content"`
	Progress  int    `json:"progress"`
	UserID    uint   `json:"user_id"`
}
