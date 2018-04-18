package main

import (
	"errors"
	"log"

	"../../api/src/routes/user"
	"../../lib"
	"github.com/jmoiron/sqlx"
)

func getUserIDFromUsername(db *sqlx.DB, senderUsername, receiverUsername string) (string, string, error) {
	var users []lib.User
	err := db.Select(&users, "SELECT id, username FROM Users WHERE username in ($1, $2)", senderUsername, receiverUsername)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect users data in database" + err.Error()))
		return "", "", errors.New("Failed to gather users data in the database")
	}
	var senderID, receiverID string
	for _, user := range users {
		if user.Username == senderUsername {
			senderID = user.ID
		} else if user.Username == receiverUsername {
			receiverID = user.ID
		}
	}
	if senderID == "" {
		return "", "", errors.New("SenderID doesn't exists in the database - can't be empty")
	}
	if receiverID == "" {
		return "", "", errors.New("ReceiverID doesn't exists in the database - can't be empty")
	}
	return senderID, receiverID, nil
}

func insertMessage(db *sqlx.DB, senderID, receiverID, content string) error {
	stmt, err := db.Preparex(`INSERT INTO Messages (senderID, receiverID, content) VALUES ($1, $2, $3)`)
	defer stmt.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert message " + err.Error()))
		return err
	}
	rows, err := stmt.Queryx(senderID, receiverID, content)
	rows.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert message " + err.Error()))
		return err
	}
	return nil
}

func messageInDB(db *sqlx.DB, senderUsername, receiverUsername, content string) error {
	senderID, receiverID, err := getUserIDFromUsername(db, senderUsername, receiverUsername)
	if err != nil {
		return err
	}
	err = insertMessage(db, senderID, receiverID, content)
	if err != nil {
		return err
	}
	errCode, errContent := user.PushNotif(db, "message", senderID, receiverID, true)
	if errCode != 0 || errContent != "" {
		return errors.New(errContent)
	}
	return nil
}
