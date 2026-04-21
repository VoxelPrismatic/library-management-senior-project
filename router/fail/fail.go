package fail

import (
	"net/http"
	"voxelprismatic/library-management-senior-project/db"

	"github.com/a-h/templ"
)

// Attempt to render the element, or return the error if it failed
func Render(p *RoutingParams, elem templ.Component) {
	err := elem.Render(p.Req.Context(), p.W)
	if err != nil {
		http.Error(p.W, err.Error(), http.StatusInternalServerError)
	}
}

// Returns 'true' if the form failed to parse
// Usage: `if fail.Form(p) { return }`
func Form(p *RoutingParams) bool {
	err := p.Req.ParseForm()
	if err != nil {
		http.Error(p.W, err.Error(), http.StatusBadRequest)
		return true
	}
	return false
}

// Returns 'true' if the user does NOT meet the minimum role requirements
// Usage: `if fail.Auth(p, UserRoleLibrarian) { return }`
func Auth(p *RoutingParams, minLevel db.UserRoleFlag) bool {
	if minLevel <= p.User.Roles {
		return false
	}

	p.W.Header().Set("X-Auth-Missing", minLevel.String())
	p.W.Header().Set("X-Auth-Current", p.User.Roles.String())
	http.Error(p.W, "Forbidden", http.StatusForbidden)
	return true
}

func Redirect(p *RoutingParams) {
	p.W.Header().Set("X-Redirect-Reason", "404: "+p.SubPathTree(true))
	http.Redirect(p.W, p.Req, p.SubPathTree(false), http.StatusPermanentRedirect)
}
