package common

import (
	"database/sql/driver"
	"time"
)

type CustomTime time.Time

func (this *CustomTime) Scan(v interface{}) error {
	if v == nil {
		*this = CustomTime(time.Time{})
		return nil
	}
	t, ok := v.(time.Time)
	if ok {
		*this = CustomTime(t)
		return nil
	}

	s, ok := v.([]uint8)
	if ok {
		parse, err := time.Parse("2006-01-02 15:04:05", string(s))
		if err != nil {
			return err
		}
		*this = CustomTime(parse)
	}
	return nil
}

func (this CustomTime) Value() (driver.Value, error) {
	return time.Time(this).Local().Format("2006-01-02 15:04:05"), nil
}
