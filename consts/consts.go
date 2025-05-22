package consts

const (
	RoleIdAdmin    = iota + 1
	RoleIdManager  = 2
	RoleIdAttendee = 3

	RoleAdmin    = "ADMIN"
	RoleManager  = "MANAGER"
	RoleAttendee = "ATTENDEE"

	DefaultPageSize = 10
	DefaultPage     = 1

	PermissionUserCreate = "user.create" // Permission to create a new user
	PermissionUserUpdate = "user.update" // Permission to update an existing user's information
	PermissionUserFetch  = "user.fetch"  // Permission to fetch a specific user's data
	PermissionUserList   = "user.list"   // Permission to list all users
	PermissionUserDelete = "user.delete" // Permission to delete a user
)

var RoleMap = map[int]string{
	RoleIdAdmin:    RoleAdmin,
	RoleIdManager:  RoleManager,
	RoleIdAttendee: RoleAttendee,
}
