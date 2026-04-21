package db

import (
	"net/http"
)

func CookieAuth(_ http.ResponseWriter, _ *http.Request) User {
	// TO-DO: Implement cookie authentication with JWT
	return User{
		Roles: UserRoleAdmin,
	}
}
