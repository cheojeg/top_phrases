package util

import (
	"database/sql"
	"time"
)

func ConvertTimeToNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
