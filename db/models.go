package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Returns how many copies are currently available to borrow
func (b *BookWork) AvailableCopies() int {
	var activeCount int64
	db.Model(&Checkout{}).
		Where("book_isbn = ? AND returned_at = ?", b.ISBN, NilTime).
		Count(&activeCount)

	available := b.TotalCopies - int(activeCount)
	if available < 0 {
		return 0
	}
	return available
}

func (c *Loan) IsActive() bool {
	return c.ReturnedAt.IsZero()
}

func (c *Loan) IsOverdue() bool {
	return c.IsActive() && time.Now().After(c.DueAt)
}

// Marks a book as returned
func (c *Loan) Return() error {
	c.ReturnedAt = time.Now()
	return Save(c)
}

func (c *Loan) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
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
	Book        Book `gorm:"foreignKey:BookISBN;references:ISBN"`
	User        User `gorm:"foreignKey:UserID;references:ID"`
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
