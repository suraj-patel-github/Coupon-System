package utils

import "time"

func ParseTimestamp(t string) (time.Time, error) {
    return time.Parse("2006-01-02", t)
}
