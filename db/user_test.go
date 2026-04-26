package db

import (
	"testing"
	"time"
)

func TestUserCheckedOut(t *testing.T) {
	tx := TestDb()
	defer tx.Rollback()

	user := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}
	tx.Save(&user)

	book := BookWork{
		ID:    "test-book-id",
		Title: "Test Book",
	}
	tx.Save(&book)

	copy := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	tx.Save(&copy)

	loan := Loan{
		BookCopyID:   copy.ID,
		UserID:       user.ID,
		DateCheckout: time.Now().Add(-DAY),
		DateReturned: NilTime,
	}
	tx.Save(&loan)

	checkedOut, err := user.CheckedOut()
	if err != nil {
		t.Fatalf("user checked out: unexpected error: %v", err)
	}
	if len(checkedOut) != 1 {
		t.Fatalf("user checked out: expected 1 loan, got %d", len(checkedOut))
	}
}

func TestUserHasOverdueBooks(t *testing.T) {
	tx := TestDb()
	defer tx.Rollback()

	user := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}
	tx.Save(&user)

	book := BookWork{
		ID:    "test-book-id",
		Title: "Test Book",
	}
	tx.Save(&book)

	copy := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	tx.Save(&copy)

	loan := Loan{
		BookCopyID:   copy.ID,
		UserID:       user.ID,
		DateCheckout: time.Now().Add(-LOAN_DURATION * 2),
		DateReturned: NilTime,
	}
	tx.Save(&loan)

	hasOverdue, err := user.HasOverdueBooks()
	if err != nil {
		t.Fatalf("user overdue books: unexpected error: %v", err)
	}
	if !hasOverdue {
		t.Fatalf("user overdue books: expected true, got false")
	}
}
