package db

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SqlStringList []string

func (bstr *SqlStringList) Scan(value any) error {
	if value == nil {
		*bstr = nil
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to string", value, value)
	}

	*bstr = strings.Split(str, "\x98") // 0x98 = Start of String
	return nil
}

func (bstr SqlStringList) Value() (driver.Value, error) {
	if len(bstr) == 0 {
		return nil, nil
	}
	return strings.Join(bstr, "\x98"), nil // 0x98 = Start of String
}

type SqlUUID struct {
	uuid.UUID
}

func (buuid *SqlUUID) Scan(value any) error {
	if value == nil {
		*buuid = SqlUUID{}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to string", value, value)
	}

	obj, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*buuid = SqlUUID{obj}

	return nil
}

func (buuid SqlUUID) Value() (driver.Value, error) {
	return buuid.String(), nil
}

func (buuid SqlUUID) IsEmpty() bool {
	for _, b := range buuid.UUID {
		if b != 0 {
			return false
		}
	}
	return true
}

var NilTime = time.Time{}

type BaseModel struct {
	gorm.Model
	ID SqlUUID `gorm:"type:text;primaryKey"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID.IsEmpty() {
		u := uuid.New()
		m.ID = SqlUUID{u}
	}
	return nil
}
