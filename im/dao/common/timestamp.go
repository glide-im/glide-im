package common

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Timestamp time.Time

func (u Timestamp) Value() (driver.Value, error) {
	tTime := time.Time(u)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (u *Timestamp) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*u = Timestamp(vt)
	default:
		return errors.New(fmt.Sprintf("%v is not type time.Time", v))
	}
	return nil
}

func (u Timestamp) MarshalJSON() ([]byte, error) {
	s := strconv.FormatInt(time.Time(u).Unix(), 10)
	return []byte(s), nil
}

func (u *Timestamp) UnmarshalJSON(bytes []byte) error {
	r, err := strconv.ParseInt(string(bytes), 10, 64)
	*u = Timestamp(time.Unix(r, 0))
	return err
}

func (u *Timestamp) String() string {
	return strconv.FormatInt(time.Time(*u).Unix(), 10)
}

func (u *Timestamp) Unix() int64 {
	return time.Time(*u).Unix()
}

func NowTimestamp() Timestamp {
	return Timestamp(time.Now())
}
