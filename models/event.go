package models

import "time"

type Event struct {
	ID          int        `json:"id" gorm:"column:id"`
	Title       string     `json:"title" gorm:"column:title"`
	Description *string    `json:"description" gorm:"column:description"`
	Location    *string    `json:"location" gorm:"column:location"`
	StartTime   *time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime     *time.Time `json:"end_time" gorm:"column:end_time"`
	IsPublic    *bool      `json:"is_public" gorm:"column:is_public"`
	Limit       *int       `json:"limit" gorm:"column:limit"`
	CreatedBy   int        `json:"created_by" gorm:"column:created_by"`
	CreatedAt   time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt   time.Time  `json:"-" gorm:"column:updated_at"`
	Attendees   []User     `json:"attendee" gorm:"many2many:event_attendees;"`
}
type EventAttendee struct {
	EventID int   `json:"event_id" gorm:"column:event_id"`
	UserID  int   `json:"user_id" gorm:"column:user_id"`
	RSVP    bool  `json:"rsvp" gorm:"column:rsvp"`
	Event   Event `gorm:"foreignKey:ID;references:EventID"`
	User    User  `gorm:"foreignKey:ID;references:UserID"`
}
