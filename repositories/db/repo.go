package db

import (
	"gorm.io/gorm"
)

type Repository struct {
	client *gorm.DB
}

func NewRepository(client *gorm.DB) *Repository {
	return &Repository{
		client: client,
	}
}
