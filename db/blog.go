package db

import (
	"time"

	"gorm.io/gorm"
)

var _ = Migrate(Blog{}, BlogEntry{})

type Blog struct {
	BaseModel
	Title  string
	Body   string
	Date   time.Time
	Tags   SqlStringList `gorm:"type:text"`
	User   User
	UserID SqlUUID `gorm:"type:text"`
}

type BlogEntry struct {
	gorm.Model
	Blog        Blog
	BlogID      SqlUUID `gorm:"primaryKey;type:text"`
	BookWork    BookWork
	BookWorkID  SqlUUID `gorm:"primaryKey;type:text"`
	Rank        int
	Description string
}

func (b Blog) Entries() ([]BlogEntry, error) {
	db := Db()
	ret := []BlogEntry{}
	status := db.Where(&BlogEntry{BlogID: b.ID}).
		Order("rank ASC").
		Find(&ret)
	return ret, status.Error
}
