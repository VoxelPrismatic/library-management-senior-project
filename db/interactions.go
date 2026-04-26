package db

import (
	"time"

	"gorm.io/gorm"
)

var _ = Migrate(UserLikes{}, UserBookmarks{}, UserRatings{})

type UserLikes struct {
	gorm.Model
	User   User
	UserID SqlUUID `gorm:"primaryKey;type:text"`
	Blog   Blog
	BlogID SqlUUID `gorm:"primaryKey;type:text"`
}

type UserBookmarks struct {
	gorm.Model
	User       User
	UserID     SqlUUID `gorm:"primaryKey;type:text"`
	BookWork   BookWork
	BookWorkID SqlUUID `gorm:"primaryKey;type:text"`
	Date       time.Time
}

type UserRatings struct {
	gorm.Model
	User       User
	UserID     SqlUUID `gorm:"primaryKey;type:text"`
	BookWork   BookWork
	BookWorkID SqlUUID `gorm:"primaryKey;type:text"`
	Rating     int
	Comment    string
}
