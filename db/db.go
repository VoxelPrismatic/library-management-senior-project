package db

import (
	"fmt"
	"reflect"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Create and connect to the database.
func connect() *gorm.DB {
	var target gorm.Dialector
	if testing.Testing() {
		target = sqlite.Open("file::memory:?cache=shared")
	} else {
		target = sqlite.Open("senior-library.db")
	}
	db, err := gorm.Open(target, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Logger.LogMode(logger.Info)

	return db
}

var _db = connect()
var _tx *gorm.DB = nil

// Automatically migrate several structs to the database
func Migrate(models ...any) bool {
	for _, model := range models {
		err := _db.AutoMigrate(model)
		if err != nil {
			fmt.Printf("\x1b[91;1mpanic: %s\x1b[0m\n", reflect.TypeOf(model).Name())
			panic(err)
		}
	}

	return true
}

// Retrieve a copy of the database pointer.
func Db() *gorm.DB {
	if testing.Testing() {
		return _tx
	}
	return _db
}

func MustSave(obj any) {
	state := _db.Save(obj)
	if state.Error != nil {
		panic(state.Error)
	}
}

func TestDb() *gorm.DB {
	if !testing.Testing() {
		panic("should never be called outside of testing environment")
	}
	_tx = _db.Begin()
	return _tx
}
