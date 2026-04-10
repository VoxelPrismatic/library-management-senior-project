package user

import (
	"voxelprismatic/library-management-senior-project/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = db.Migrate(User{})

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
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

type UserStatusFlag int

const (
	UserStatusActive  UserStatusFlag = 1 << iota
	UserStatusLocked                 // TO-DO: Check if user has outstanding fees and remove this redundant lock
	UserStatusDeleted                // For audit purposes; we may choose to anonymize any data
)
