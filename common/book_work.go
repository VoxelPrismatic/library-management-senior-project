package common

import (
	"time"

	"gorm.io/gorm"
)

type BookWork struct {
	gorm.Model
	Isbn_13       int `gorm:"primaryKey"`
	Isbn_10       int
	GoogleBooksID string
	OpenLibraryID string
	WorldCatID    string
	Title         string
	Authors       SqlStringList
	Tags          SqlStringList
	Published     time.Time
	Language      string
	PageCount     int
	Description   string
}
