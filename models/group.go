package models

import "gorm.io/gorm"

// Group represents a shared expense group (e.g., roommates, trip)
type Group struct {
	gorm.Model
	Name      string        `json:"name" gorm:"not null"`
	CreatedBy uint          `json:"created_by"`
	Members   []GroupMember `json:"members,omitempty" gorm:"foreignKey:GroupID"`
}

// GroupMember is the many-to-many join table between Group and User
type GroupMember struct {
	gorm.Model
	GroupID uint `json:"group_id" gorm:"not null"`
	UserID  uint `json:"user_id" gorm:"not null"`
}
