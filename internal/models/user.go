package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primarykey"`
	Username  string `gorm:"size:50;not null;unique"`
	Password  string `gorm:"size:255;not null"`
	Email     string `gorm:"size:100;not null;unique"`
	Name      string `gorm:"size:50;not null"`
	Status    int    `gorm:"default:1"`
	Roles     []Role `gorm:"many2many:user_roles;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
