package lib

import (
	"time"
)

// User is the data structure of the table User from PostgreSQL
type User struct {
	ID                     string     `db:"id"`
	Username               string     `db:"username"`
	Email                  string     `db:"email"`
	Lastname               string     `db:"lastname"`
	Firstname              string     `db:"firstname"`
	Password               string     `db:"password"`
	CreatedAt              time.Time  `db:"created_at"`
	RandomToken            string     `db:"random_token"`
	PictureURL_1           string     `db:"picture_url_1"`
	PictureURL_2           string     `db:"picture_url_2"`
	PictureURL_3           string     `db:"picture_url_3"`
	PictureURL_4           string     `db:"picture_url_4"`
	PictureURL_5           string     `db:"picture_url_5"`
	Biography              string     `db:"biography"`
	Birthday               *time.Time `db:"birthday"`
	Genre                  string     `db:"genre"`
	InterestingIn          string     `db:"interesting_in"`
	Latitude               *float64   `db:"latitude"`
	Longitude              *float64   `db:"longitude"`
	GeolocalisationAllowed bool       `db:"geolocalisation_allowed"`
}
