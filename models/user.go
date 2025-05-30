package models

import "time"

type (
	User struct {
		ID        int
		Email     string
		Password  string `json:"-"`
		FirstName string
		LastName  string
		RoleID    int
		CreatedAt time.Time
		UpdatedAt time.Time
		Events    []Event `gorm:"many2many:event_attendees;"`
	}

	RolePermission struct {
		RoleID       int
		PermissionID int
	}

	Permission struct {
		ID         int
		Permission string
	}
)
