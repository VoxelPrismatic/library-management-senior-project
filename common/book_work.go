package common

import (
	"time"

	"gorm.io/gorm"
)

type BookWork struct {
	gorm.Model
	ID string `gorm:"primaryKey"` // Google Books Volume ID

	Title         string
	Subtitle      string
	Authors       SqlStringList
	Publisher     string
	PublishedDate time.Time
	Version       string

	Isbn13 int64
	Isbn10 int64

	Description string
	PageCount   int
	IsMature    bool
	Categories  SqlStringList

	CoverThumb string
	CoverImage string
}
