package models

import "time"

type (
	User struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"`
		Password  string    `json:"-"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		RoleID    int       `json:"-"`
		CreatedAt time.Time `json:"-"`
		UpdatedAt time.Time `json:"-"`
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
