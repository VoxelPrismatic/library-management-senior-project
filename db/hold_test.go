package db

import (
	"testing"
	"time"
)

func TestHoldStatusCompleted(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	hold := Hold{
		FulfilledDate: time.Now(),
	}

	status, err := hold.Status()
	if err != nil {
		t.Fatalf("hold status: unexpected error: %v", err)
	}
	if status != HoldCompleted {
		t.Fatalf("hold status: expected HoldCompleted, got %d", status)
	}
}

func TestHoldStatusCancelled(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	hold := Hold{
		CancelledDate: time.Now(),
	}

	status, err := hold.Status()
	if err != nil {
		t.Fatalf("hold status: unexpected error: %v", err)
	}
	if status != HoldCancelled {
		t.Fatalf("hold status: expected HoldCancelled, got %d", status)
	}
}

func TestHoldStatusQueued(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	user := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Status:    UserStatusActive,
	}
	db.Save(&user)

	book := BookWork{
		ID:    "test-book-id",
		Title: "Test Book",
	}
	db.Save(&book)

	hold := Hold{
		User:       user,
		UserID:     user.ID,
		BookWorkID: SqlUUID{},
	}

	status, err := hold.Status()
	if err != nil {
		t.Fatalf("hold status: unexpected error: %v", err)
	}
	if status != HoldQueued {
		t.Fatalf("hold status: expected HoldQueued, got %d", status)
	}
}
