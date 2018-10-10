package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/satori/go.uuid"
)

type key int

const UserContext key = 0

func GetSessionToken(r *http.Request) string {
	sessionToken := r.Header.Get("X-Session-Token")
	if sessionToken == "" {
		mlog.Error("No auth token sent in request header")
		// Token not sent in request
		return ""
	}

	return sessionToken
}

func CheckLoginStatus(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := GetSessionToken(r)
		if token != "" {
			u, loggedIn := GetUserBySessionToken(token)
			if loggedIn {
				context.Set(r, UserContext, u)
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

// ClearSession deletes the a session token from the database, effectively logging a user out
func ClearSession(sessionToken string) {
	DeleteSessionToken(sessionToken)
}

// SessionsHandlerLogin is a POST API endpoint that takes a username and password, and attempts to log in a user
// The user will be redirected to "/" if the credentials are incorrect, or will be forwarded to "/tasks" if successful
func SessionsHandlerLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		// Return bad login message
		return
	}
	if username != "" && password != "" {
		u, err := GetUserByUsername(username)
		if err != nil {
			mlog.Error(err)
			// Return bad login message
			return
		}
		if u.CheckPasswordHash(password) {
			sessionToken, err := uuid.NewV4()
			if err != nil {
				mlog.Error(err)
				// Return bad login message
				return
			}
			CreateSession(sessionToken.String(), u)
			mlog.Info(u.Username, "has successfully logged in")
			// TODO: Return sessionToken as JSON
			fmt.Fprint(w, sessionToken)
		}
	}
}

// SessionsHandlerLogout is a POST API endpoint that logs a user out, and deletes the associated session token
func SessionsHandlerLogout(w http.ResponseWriter, r *http.Request) {
	token := GetSessionToken(r)
	if token != "" {
		_, loggedIn := GetUserBySessionToken(token)
		if loggedIn {
			ClearSession(token)
			// Return successfully logged out notification
		}
	} else {
		// Return bad login message
	}
}

func SessionsHandlerCheck(w http.ResponseWriter, r *http.Request) {
	token := GetSessionToken(r)
	if token != "" {
		_, loggedIn := GetUserBySessionToken(token)
		if loggedIn {
			http.Error(w, "OK", http.StatusOK)
			return
		} else {
			ClearSession(token)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}
