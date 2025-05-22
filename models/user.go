package models

import "time"

type (
	User struct {
		ID        int
		Email     string
		Password  string
		FirstName string
		LastName  string
		RoleID    int
		CreatedAt time.Time
		UpdatedAt time.Time
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
