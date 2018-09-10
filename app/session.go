package app

import (
	"time"

	"github.com/lib/pq"
)

// duration in hours before a session expires
// Needs to be moved to a configuration file
var sessionDuration = 24

var userSessionTableStatement = `CREATE TABLE IF NOT EXISTS usersessions (
	sessiontoken VARCHAR(36) PRIMARY KEY,
	userid INT NOT NULL,
	logintime DATETIME NOT NULL,
	lastseentime DATETIME NOT NULL
)`

// UserSession is the basic structure of a user's session, and forms the basis for storing that data in the database
// A UserSession should be tracked and verified as a user moves through the web application
type UserSession struct {
	SessionToken string
	UserID       int
	LoginTime    pq.NullTime
	LastSeenTime pq.NullTime
}

// CreateSession creates a session and stores it in the database for a given user, with a session token
func CreateSession(sessionToken string, user *User) {
	stmt, err := store.db.Prepare("INSERT INTO usersessions(sessiontoken, userid, logintime, lastseentime) VALUES(?, ?, ?, ?)")
	_, err = stmt.Exec(sessionToken, user.ID, time.Now(), time.Now())
	if err != nil {
		mlog.Error(err)
		return
	}
}

// GetUserBySessionToken takes a sessionToken, and determines if it is a valid token
// If the token is valid, it will verify if the user is currently logged in based on sessionDuration and LoginTime
// If the user is logged in, it will update the UserSession's LastSeenTime
// Returns nil and false if the session is not valid, or Returns the user and true if the user is logged in
func GetUserBySessionToken(sessionToken string) (*User, bool) {
	var session UserSession

	err := store.db.QueryRow("SELECT sessiontoken, userid, logintime, lastseentime FROM usersessions WHERE sessiontoken = ?", sessionToken).Scan(&session.SessionToken, &session.UserID, &session.LoginTime, &session.LastSeenTime)
	if err != nil {
		mlog.Error(err)
		return nil, false
	}

	// get the time of number of hours of sessionDuration ago (eg. 24 hours ago)
	// then compare it to the LoginTime, and see if the user was Logged in after sessionDurationAgo
	sessionDurationAgo := time.Now().Add(time.Duration(-1*sessionDuration) * time.Hour)
	loggedIn := session.LoginTime.Time.After(sessionDurationAgo)

	// Update or Delete the session token if they are not logged in anymore
	if loggedIn {
		_, err := store.db.Exec("UPDATE usersessions SET lastseentime = ? WHERE sessiontoken = ?", time.Now(), sessionToken)
		if err != nil {
			mlog.Error(err)
			return nil, false
		}
	} else {
		DeleteSessionToken(session.SessionToken)
	}
	user, err := GetUserByID(session.UserID)
	return user, true
}

// GetSessionsByUser returns all the sessions in the database for a given userid
func GetSessionsByUser(userID int) []UserSession {
	var sessions []UserSession

	rows, err := store.db.Query("SELECT sessiontoken, userid, logintime, lastseentime FROM usersessions WHERE userid = ?", userID)
	if err != nil {
		mlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var session UserSession
		err := rows.Scan(&session.SessionToken, &session.UserID, &session.LoginTime, &session.LastSeenTime)
		if err != nil {
			mlog.Error(err)
			return sessions
		}
		sessions = append(sessions, session)
	}
	return sessions
}

// DeleteSessionToken removes a session from the database
// This should be called when a token is no longer valid
func DeleteSessionToken(sessionToken string) {
	stmt, err := store.db.Prepare("DELETE FROM usersessions WHERE sessiontoken = ?")
	_, err = stmt.Exec(sessionToken)
	if err != nil {
		mlog.Error(err)
		return
	}
}
