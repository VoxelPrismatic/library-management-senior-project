package db

import (
	"time"
)

var _ = Migrate(Hold{})

// When a user wants to get in line to check out a book
type Hold struct {
	BaseModel
	BookWork      BookWork
	BookWorkID    SqlUUID `gorm:"type:text"`
	User          User
	UserID        SqlUUID `gorm:"type:text"`
	RequestedDate time.Time
	FulfilledDate time.Time
	CancelledDate time.Time
}

func (h Hold) Status() (HoldStatusFlag, error) {
	if !h.FulfilledDate.IsZero() {
		return HoldCompleted, nil
	}

	if !h.CancelledDate.IsZero() {
		return HoldCancelled, nil
	}

	u := h.GetUser()
	if u.Status == UserStatusDeleted {
		return HoldRevoked, nil
	}

	if state, err := u.HasOverdueBooks(); err != nil {
		return 0, err
	} else if state {
		return HoldPostponed, nil
	}

	if count, err := u.CheckedOut(); err != nil {
		return 0, err
	} else if len(count) >= LOAN_LIMIT {
		return HoldPostponed, nil
	}

	return HoldQueued, nil
}

func (h Hold) GetUser() User {
	db := Db()
	if h.User.ID.IsEmpty() {
		db.Where(&User{BaseModel: BaseModel{ID: h.UserID}}).First(&h.User)
	}
	return h.User
}
