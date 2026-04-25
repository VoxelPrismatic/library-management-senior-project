package db

import (
	"testing"
)

func TestAvailableCopiesTotal(t *testing.T) {
	setupTestDb(t)

	book := BookWork{
		ID:    "test-book-id",
		Title: "Test Book",
	}
	db.Save(&book)

	copy1 := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	copy2 := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	db.Save(&copy1)
	db.Save(&copy2)

	counts, err := book.AvailableCopies(true)
	if err != nil {
		t.Fatalf("available copies: unexpected error: %v", err)
	}

	result, ok := counts[BookFmtPaperback]
	if !ok {
		t.Fatalf("available copies: expected BookFmtPaperback in results")
	}
	if result.Total != 2 {
		t.Fatalf("available copies: expected total 2, got %d", result.Total)
	}
	if result.Available != 2 {
		t.Fatalf("available copies: expected available 2, got %d", result.Available)
	}
}

func TestAvailableCopiesWithLoan(t *testing.T) {
	setupTestDb(t)

	user := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}
	db.Save(&user)

	book := BookWork{
		ID:    "test-book-id-2",
		Title: "Test Book 2",
	}
	db.Save(&book)

	copy1 := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	copy2 := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	db.Save(&copy1)
	db.Save(&copy2)

	loan := Loan{
		BookCopyID:   copy1.ID,
		UserID:       user.ID,
		DateCheckout: NilTime,
		DateReturned: NilTime,
	}
	db.Save(&loan)

	counts, err := book.AvailableCopies(true)
	if err != nil {
		t.Fatalf("available copies with loan: unexpected error: %v", err)
	}

	result, ok := counts[BookFmtPaperback]
	if !ok {
		t.Fatalf("available copies with loan: expected BookFmtPaperback in results")
	}
	if result.Total != 2 {
		t.Fatalf("available copies with loan: expected total 2, got %d", result.Total)
	}
	if result.Available != 1 {
		t.Fatalf("available copies with loan: expected available 1, got %d", result.Available)
	}
}
