package db

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

var _ = Migrate(User{})

type User struct {
	BaseModel
	Roles     UserRoleFlag
	FirstName string
	LastName  string
	Email     string `gorm:"unique;not null"`
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

func (u *User) SetFirstName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("first name cannot be blank")
	}
	if len(name) > MAX_NAME_LEN {
		return fmt.Errorf("first name cannot exceed %d characters", MAX_NAME_LEN)
	}

	u.FirstName = name
	return nil

}

func (u *User) SetLastName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("last name cannot be blank")
	}

	if len(name) > MAX_NAME_LEN {
		return fmt.Errorf("last name cannot exceed %d characters", MAX_NAME_LEN)
	}

	u.LastName = name
	return nil
}

func (u *User) SetEmail(addr string) error {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return fmt.Errorf("email cannot be blank")
	}

	if len(addr) > MAX_EMAIL_LEN {
		return fmt.Errorf("email cannot exceed %d characters", MAX_EMAIL_LEN)
	}

	if !emailRegexp.MatchString(addr) {
		return fmt.Errorf("email is malformatted")
	}

	addrLower := strings.ToLower(addr)
	lookup := User{}
	db.Where(User{Email: addrLower}).First(&lookup)
	if lookup.Email != "" && lookup.ID.UUID != u.ID.UUID {
		return fmt.Errorf("email already in use")
	}

	u.Email = addrLower
	return nil
}

func TestSecretStrength(secret string) error {
	if len(secret) < MIN_SECRET_LEN {
		return fmt.Errorf("secret too short")
	}

	if len(secret) > MAX_SECRET_LEN {
		return fmt.Errorf("secret cannot exceed %d characters", MAX_SECRET_LEN)
	}

	inclUpper := false
	inclLower := false
	inclDigit := false
	inclSymbol := false

	for _, r := range secret {
		switch {
		case unicode.IsUpper(r):
			inclUpper = true
		case unicode.IsLower(r):
			inclLower = true
		case unicode.IsDigit(r):
			inclDigit = true
		case unicode.IsSymbol(r):
			inclSymbol = true
		case unicode.IsPunct(r):
			inclSymbol = true
		}
	}

	missing := []string{}
	if !inclUpper {
		missing = append(missing, "an uppercase letter")
	}
	if !inclLower {
		missing = append(missing, "a lowercase letter")
	}
	if !inclDigit {
		missing = append(missing, "a number")
	}
	if !inclSymbol {
		missing = append(missing, "a symbol")
	}

	if len(missing) > 1 {
		return fmt.Errorf("secret is missing %s, and %s", strings.Join(missing[:len(missing)-1], ","), missing[len(missing)-1])
	} else if len(missing) == 1 {
		return fmt.Errorf("secret is missing %s", missing[0])
	}

	return nil
}

func (u *User) TestSecret(secret string) bool {
	return u.HashSecret(secret) == u.Secret
}

func (u *User) SetSecret(secret, verify string) error {
	if err := TestSecretStrength(secret); err != nil {
		return fmt.Errorf("secret has problems")
	}

	if secret != verify {
		return fmt.Errorf("secret mismatch")
	}

	u.Secret = u.HashSecret(secret)
	return nil
}

func (u *User) HashSecret(pass string) string {
	cycle := uint16(len(u.Email))
	for _, r := range pass {
		cycle <<= 1
		cycle ^= uint16(r)
	}

	salt := fmt.Appendf([]byte{}, "%s.%d.%s", u.Email, cycle, pass)

	for i := 0; i < int(cycle); i++ {
		hash := crypto.SHA512.New()
		hash.Write(salt)
		salt = hash.Sum(nil)
	}

	return hex.EncodeToString(salt)
}
