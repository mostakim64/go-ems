package consts

import "time"

const (
	RoleIdAdmin    = iota + 1
	RoleIdManager  = 2
	RoleIdAttendee = 3

	RoleAdmin    = "ADMIN"
	RoleManager  = "MANAGER"
	RoleAttendee = "ATTENDEE"

	DefaultPageSize = 10
	DefaultPage     = 1

	PermissionUserCreate             = "user.create"       // Permission to create a new user
	PermissionUserUpdate             = "user.update"       // Permission to update an existing user's information
	PermissionUserFetch              = "user.fetch"        // Permission to fetch a specific user's data
	PermissionUserList               = "user.list"         // Permission to list all users
	PermissionUserDelete             = "user.delete"       // Permission to delete a user
	PermissionListAttendee           = "user.listAttendee" // list attendee
	PermissionFetchAllUserAsAttendee = "user.fetchAllUserAsAttendee"

	PermissionEventCreate       = "event.create" // Permission to create a new event
	PermissionEventUpdate       = "event.update" // Permission to update an existing event
	PermissionEventFetch        = "event.fetch"  // Permission to fetch a specific event
	PermissionEventList         = "event.list"   // Permission to list events
	PermissionEventDelete       = "event.delete" // Permission to delete an event
	PermissionFetchAllEvent     = "event.fetchAllEvent"
	PermissionFetchOwnEvent     = "event.fetchOwnEvent"
	PermissionFetchInvitedEvent = "event.fetchInvitedEvent"

	StatusInvited  = 1
	StatusAccepted = 2
	StatusRejected = 3

	EventReminderInterval = time.Duration(10 * time.Minute)
)

var RoleMap = map[int]string{
	RoleIdAdmin:    RoleAdmin,
	RoleIdManager:  RoleManager,
	RoleIdAttendee: RoleAttendee,
}
