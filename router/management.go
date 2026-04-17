package router

import (
	"fmt"
	"net/http"
	"net/url"
	"voxelprismatic/library-management-senior-project/router/fail"
	"voxelprismatic/library-management-senior-project/web/pages"
	"voxelprismatic/library-management-senior-project/web/user"
)

func ManagementRouter(p *fail.RoutingParams) {
	if fail.Auth(p, user.UserRoleLibrarian) {
		return
	}

	switch p.Pop() {
	case "books":
		ManagementBooksRouter(p)
	default:
		fail.Redirect(p)
	}
}

func ManagementBooksRouter(p *fail.RoutingParams) {
	if fail.Auth(p, user.UserRoleLibrarian) {
		return
	}

	switch p.Pop() {
	case "add":
		HandleManagementBooksAdd(p)
	default:
		fail.Redirect(p)
	}
}

func HandleManagementBooksAdd(p *fail.RoutingParams) {
	switch p.Req.Method {
	case http.MethodGet:
		query := p.Req.URL.Query().Get("q")
		fail.Render(p, pages.BookMgmtSearchFull(query))
	case http.MethodPost:
		if fail.Form(p) {
			return
		}
		query := p.Req.Form.Get("q")
		p.W.Header().Set(
			"hx-replace-url",
			fmt.Sprintf("%s?q=%s", p.Req.URL.String(), url.QueryEscape(query)),
		)
		fail.Render(p, pages.BookMgmtSearchGrid(query))

	}
}
