package db

import (
	"time"

	"github.com/google/uuid"
)

var _ = Migrate(User{})

type User struct {
	BaseModel
	Roles     UserRoleFlag
	FirstName string
	LastName  string
	Email     string
	Secret    string
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
	ret := []Loan{}
	status := db.Model(&Loan{}).Joins("users").Where(
		"date_returned = ?", NilTime,
	).Find(&ret)
	return ret, status.Error
}

func (u User) HasOverdueBooks() (bool, error) {
	count := int64(0)
	status := db.Model(&Loan{}).Joins("users").Where(
		"date_returned = ?", NilTime,
	).Where(
		"date_checkout < ?", time.Now().Add(-LOAN_DURATION),
	).Count(&count)
	return count == 0, status.Error
}

func (u User) Partial() UserPartial {
	return UserPartial{
		ID:        u.ID.String(),
		Roles:     u.Roles,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

type UserPartial struct {
	ID        string
	Roles     UserRoleFlag
	FirstName string
	LastName  string
	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`
}

func (p UserPartial) Fetch() (User, error) {
	id, err := uuid.Parse(p.ID)
	if err != nil {
		return User{}, err
	}

	ret := User{BaseModel: BaseModel{ID: SqlUUID{id}}}
	err = db.Where(&ret).First(&ret).Error
	return ret, err
}

func (p *UserPartial) SetTimestamp(t int64) {
	p.IssuedAt = t
	p.ExpiresAt = t + JWT_LIFETIME
}
