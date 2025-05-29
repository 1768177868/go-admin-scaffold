package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint         `gorm:"primarykey"`
	Name        string       `gorm:"size:50;not null;unique"`
	DisplayName string       `gorm:"size:100;not null"`
	Description string       `gorm:"size:255"`
	Users       []User       `gorm:"many2many:user_roles;"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
