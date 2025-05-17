package utils

import "time"

func ParseTimestamp(t string) (time.Time, error) {
    return time.Parse(time.RFC3339, t)
}
