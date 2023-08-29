package models

import (
	"time"

	"gorm.io/gorm"
)

type Alias struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string    `gorm:"primaryKey;size:255" json:"name"`
	Url      string    `gorm:"size:2038;not null" json:"url"`
	Requests []Request `gorm:"foreignKey:AliasName" json:"requests"`
}
