package models

import "time"

type BlacklistedToken struct {
	ID               uint `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time
	BlacklistedToken string `json:"blacklisted_token"`
}
