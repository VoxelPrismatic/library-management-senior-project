package db

import (
	"fmt"
	"reflect"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create and connect to the database.
func connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("senior-library.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

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

// Get the first object that matches the parameters of the filter object
// The filter object is simply a partially filled struct
// Note: This function returns a new object; it does not modify the object you pass
func GetFirst[T any](filter T) T {
	ret := new(T)
	db.Where(&filter).First(ret)
	return *ret
}

// Get many rows of data, with an optional sort.
// If sort is an empty string, no sort is performed
func GetMany[T any](filter T, sort string) []T {
	ret := new([]T)
	if sort != "" {
		db.Where(&filter).Order(sort).Find(ret)
	} else {
		db.Where(&filter).Find(ret)
	}
	return *ret
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
