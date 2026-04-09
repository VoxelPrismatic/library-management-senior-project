package common

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type SqlStringList []string

func (bstr *SqlStringList) Scan(value interface{}) error {
	if value == nil {
		*bstr = nil
		return nil
	}

	str, okay := value.(string)
	if !okay {
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
