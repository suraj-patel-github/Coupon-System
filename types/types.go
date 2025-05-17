package types

import (
	"strings"
	"time"
)

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}
