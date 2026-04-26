package db

import (
	"fmt"
	"reflect"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create and connect to the database.
func connect() *gorm.DB {
	var target gorm.Dialector
	if testing.Testing() {
		// target = sqlite.Open("file::memory:?cache=shared")
		target = sqlite.Open("testing.db")
	} else {
		target = sqlite.Open("senior-library.db")
	}
	db, err := gorm.Open(target, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

var db = connect()

// Automatically migrate several structs to the database
func Migrate(models ...any) bool {
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			fmt.Printf("\x1b[91;1mpanic: %s\x1b[0m\n", reflect.TypeOf(model).Name())
			panic(err)
		}
	}

	return true
}

// Retrieve a copy of the database pointer.
func Db() *gorm.DB {
	return db
}

func MustSave(obj any) {
	state := db.Save(obj)
	if state.Error != nil {
		panic(state.Error)
	}
}
