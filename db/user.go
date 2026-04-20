package db

var _ = Migrate(User{})

type User struct {
	BaseModel
	Roles     UserRoleFlag
	FirstName string
	LastName  string
	Email     string
	Status    UserStatusFlag
}

type UserRoleFlag int

const (
	UserRolePublic UserRoleFlag = 1 << iota
	UserRoleLibrarian
	UserRoleAdmin
)

func (f UserRoleFlag) String() string {
	switch f {
	case UserRolePublic:
		return "public"
	case UserRoleLibrarian:
		return "librarian"
	case UserRoleAdmin:
		return "administrator"
	default:
		return "undefined"
	}
}

type UserStatusFlag int

const (
	UserStatusActive  UserStatusFlag = 1 << iota
	UserStatusLocked                 // TO-DO: Check if user has outstanding fees and remove this redundant lock
	UserStatusDeleted                // For audit purposes; we may choose to anonymize any data
)
