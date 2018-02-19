package lib

import (
	"time"
)

// User is the data structure of the table User from PostgreSQL
type User struct {
	ID          string    `db:"id"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	Lastname    string    `db:"lastname"`
	Firstname   string    `db:"firstname"`
	Password    string    `db:"password"`
	CreatedAt   time.Time `db:"created_at"`
	LoginToken  string    `db:"login_token"`
	RandomToken string    `db:"random_token"`
}
