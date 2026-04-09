package common

import (
	"time"

	"gorm.io/gorm"
)

type BookWork struct {
	gorm.Model
	ID string `gorm:"primaryKey"` // Google Books Volume ID

	Title         string
	Authors       SqlStringList
	Publisher     string
	PublishedDate time.Time
	Version       string

	Isbn_13 int
	Isbn_10 int

	Description string
	PageCount   int
	IsMature    bool
	Categories  SqlStringList
}
