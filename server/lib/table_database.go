package lib

import (
	"time"
)

// User is the data structure of the table User from PostgreSQL
type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
}
