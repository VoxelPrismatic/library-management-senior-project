package db

import (
	"fmt"
	"reflect"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Create and connect to the database.
func connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("senior-library.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Logger.LogMode(logger.Info)

	return db
}

var db *gorm.DB = connect()

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
	fmt.Println("\x1b[94;1mMustSave\x1b[0m")
	state := db.Save(obj)
	if state.Error != nil {
		panic(state.Error)
	}
}
