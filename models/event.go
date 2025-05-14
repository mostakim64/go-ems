package models

import "time"

type Event struct {
	ID          int        `json:"id" gorm:"column:id"`
	Title       string     `json:"title" gorm:"column:title"`
	Description *string    `json:"description" gorm:"column:description"`
	Location    *string    `json:"location" gorm:"column:location"`
	StartTime   *time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime     *time.Time `json:"end_time" gorm:"column:end_time"`
	CreatedBy   *string    `json:"created_by" gorm:"column:created_by"`
	CreatedAt   time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt   time.Time  `json:"-" gorm:"column:updated_at"`
}
