package db

import "time"

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
	UserStatusLimited                // User has hit loan limit
	UserStatusLocked                 // TO-DO: Check if user has outstanding fees and remove this redundant lock
	UserStatusDeleted                // For audit purposes; we may choose to anonymize any data
)

func (u User) CheckedOut() ([]Loan, error) {
	db := Db()
	ret := []Loan{}
	status := db.Model(&Loan{}).
		Where("user_id = ?", u.ID).
		Where("date_returned = ?", NilTime).
		Preload("User").
		Preload("BookCopy").
		Find(&ret)
	return ret, status.Error
}

func (u User) HasOverdueBooks() (bool, error) {
	db := Db()
	count := int64(0)
	status := db.Model(&Loan{}).
		Where("user_id = ?", u.ID).
		Where("date_returned = ?", NilTime).
		Where("date_checkout < ?", time.Now().Add(-LOAN_DURATION)).
		Preload("User").
		Preload("BookCopy").
		Count(&count)
	return count != 0, status.Error
}
