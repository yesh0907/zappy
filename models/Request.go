package models

import (
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	IP        string `gorm:"size:255;not null;" json:"ip"`
	UserAgent string `gorm:"size:255;" json:"user_agent"`
	Referer   string `gorm:"size:255;" json:"referer"`
	UserId    string `gorm:"size:32" json:"user_id"`
	AliasName string `gorm:"size:255" json:"alias_name"`
}
