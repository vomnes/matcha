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
	City                   string     `db:"city"`
	ZIP                    string     `db:"zip"`
	Country                string     `db:"country"`
	Latitude               *float64   `db:"latitude"`
	Longitude              *float64   `db:"longitude"`
	GeolocalisationAllowed bool       `db:"geolocalisation_allowed"`
	Online                 bool       `db:"online"`
	Rating                 float64    `db:"rating"`
}

// Tag is the data structure of the table Tag from PostgreSQL
type Tag struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// UserTag is the data structure of the table User_Tag from PostgreSQL
type UserTag struct {
	ID     string `db:"id"`
	UserID string `db:"userid"`
	TagID  string `db:"tagid"`
}

// Like is the data structure of the table Likes from PostgreSQL
type Like struct {
	ID          string    `db:"id"`
	UserID      string    `db:"userid"`
	LikedUserID string    `db:"liked_userid"`
	CreatedAt   time.Time `db:"created_at"`
}

// Visit is the data structure of the table Visits from PostgreSQL
type Visit struct {
	ID            string    `db:"id"`
	UserID        string    `db:"userid"`
	VisitedUserID string    `db:"visited_userid"`
	CreatedAt     time.Time `db:"created_at"`
}

// FakeReport is the data structure of the table Fake_Reports from PostgreSQL
type FakeReport struct {
	ID           string    `db:"id"`
	UserID       string    `db:"userid"`
	TargetUserID string    `db:"target_userid"`
	CreatedAt    time.Time `db:"created_at"`
}

// Message is the data structure of the table Messages from PostgreSQL
type Message struct {
	ID         string    `db:"id"`
	SenderID   string    `db:"senderid"`
	ReceiverID string    `db:"receiverid"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
	IsRead     bool      `db:"is_read"`
}

// Notification is the data structure of the table Notifications from PostgreSQL
type Notification struct {
	ID           string    `db:"id"`
	TypeID       string    `db:"typeid"`
	UserID       string    `db:"userid"`
	TargetUserID string    `db:"target_userid"`
	CreatedAt    time.Time `db:"created_at"`
	IsRead       bool      `db:"is_read"`
}

// NotificationsType is the data structure of the table Notifications_Types from PostgreSQL
type NotificationsType struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
