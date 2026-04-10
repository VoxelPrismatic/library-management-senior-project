package book

import (
	"time"

	"voxelprismatic/library-management-senior-project/db"

	"gorm.io/gorm"
)

var _ = db.Migrate(BookWork{})

// A literary work with all relevant metadata
type BookWork struct {
	gorm.Model
	ID string `gorm:"primaryKey"` // Google Books Volume ID

	Title         string           // E.g. A Woman Underground
	Subtitle      string           // E.g. A Cameron Winter Mystery
	Authors       db.SqlStringList `gorm:"type:text"`
	Publisher     string
	PublishedDate time.Time
	Version       string // As provided by Google Books

	Isbn13 int64
	Isbn10 int64

	Description string
	PageCount   int
	IsMature    bool
	Categories  db.SqlStringList `gorm:"type:text"`

	CoverThumb string
	CoverImage string
}
