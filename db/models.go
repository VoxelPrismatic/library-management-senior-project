package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// One row per book title in the catalog
type Book struct {
	ISBN        string        `gorm:"primaryKey"`
	Title       string        `gorm:"not null"`
	Author      string        `gorm:"not null"`
	Description string
	CoverURL    string
	Publisher   string
	PublishedAt string
	CallNumber  string
	Genres      SqlStringList `gorm:"type:text"`
	TotalCopies int           `gorm:"default:1"`
}

// Returns how many copies are currently available to borrow
func (b *Book) AvailableCopies() int {
	active := GetMany(Checkout{BookISBN: b.ISBN, ReturnedAt: NilTime}, "")
	available := b.TotalCopies - len(active)
	if available < 0 {
		return 0
	}
	return available
}

// One row per borrowing event — handles current, overdue, and history
// ReturnedAt is zero if the book hasn't been returned yet
type Checkout struct {
	ID           string    `gorm:"primaryKey"`
	UserID       string    `gorm:"not null;index"`
	BookISBN     string    `gorm:"not null;index"`
	CheckedOutAt time.Time `gorm:"not null"`
	DueAt        time.Time `gorm:"not null"`
	ReturnedAt   time.Time
	Book         Book      `gorm:"foreignKey:BookISBN;references:ISBN"`
	User         User      `gorm:"foreignKey:UserID;references:ID"`
}

func (c *Checkout) IsActive() bool {
	return c.ReturnedAt == NilTime
}

func (c *Checkout) IsOverdue() bool {
	return c.IsActive() && time.Now().After(c.DueAt)
}

// Marks a book as returned
func (c *Checkout) Return() error {
	c.ReturnedAt = time.Now()
	return Save(c)
}

func (c *Checkout) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}

type User struct {
	ID           string    `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Name         string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	RoleLevel    int       `gorm:"default:1"`
	CreatedAt    time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

// A reservation for a book that is currently unavailable
// FulfilledAt is zero if the hold is still waiting
type Hold struct {
	ID          string    `gorm:"primaryKey"`
	UserID      string    `gorm:"not null;index"`
	BookISBN    string    `gorm:"not null;index"`
	PlacedAt    time.Time `gorm:"not null"`
	FulfilledAt time.Time
	Book        Book      `gorm:"foreignKey:BookISBN;references:ISBN"`
	User        User      `gorm:"foreignKey:UserID;references:ID"`
}

func (h *Hold) IsPending() bool {
	return h.FulfilledAt == NilTime
}

func (h *Hold) BeforeCreate(tx *gorm.DB) error {
	if h.ID == "" {
		h.ID = uuid.NewString()
	}
	return nil
}

// Creates all tables in the database on startup
var _ = Migrate(
	&Book{},
	&User{},
	&Checkout{},
	&Hold{},
)
