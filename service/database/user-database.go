package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserIDByUsername(username string) (int, error) {
	var userID int
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return userID, nil
}

func (db *appdbimpl) UserExists(userID int) bool {

	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE id = ?", userID).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

func (db *appdbimpl) UsernameExists(username string) bool {

	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (db *appdbimpl) ChangeUsername(userID int, newUsername string) error {

	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, userID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) CreateUser(username string) (int, error) {

	userID, err := db.GetUserIDByUsername(username)
	if err != nil {
		return 0, err
	}

	// If User already exists -> Return the identifier
	if userID != 0 {
		return userID, nil
	}

	// Else create a new User and return the new identifier
	result, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
	if err != nil {
		return 0, err
	}

	newUserID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(newUserID), nil
}

func (db *appdbimpl) GetUserProfile(userID int) (UserProfile, error) {

	userProfile := UserProfile{
		UserID: userID,
	}

	var username string

	err := db.c.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return userProfile, err
	}

	userProfile.Username = username

	userProfile.Photos, err = db.GetUserPhotos(userID)
	if err != nil {
		return userProfile, err
	}

	userProfile.Following, err = db.GetUserFollowing(userID)
	if err != nil {
		return userProfile, err
	}

	userProfile.Followers, err = db.GetUserFollowers(userID)
	if err != nil {
		return userProfile, err
	}

	return userProfile, nil
}
