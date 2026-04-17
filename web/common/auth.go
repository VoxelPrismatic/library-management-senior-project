package common

import (
	"net/http"
	"voxelprismatic/library-management-senior-project/web/user"
)

func CookieAuth(_ http.ResponseWriter, _ *http.Request) user.User {
	// TO-DO: Implement cookie authentication with JWT
	return user.User{
		Roles: user.UserRoleAdmin,
	}
}
