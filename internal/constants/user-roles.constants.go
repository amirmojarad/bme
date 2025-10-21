package constants

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleGuest UserRole = "user"
)

func (u UserRole) String() string {
	return string(u)
}
