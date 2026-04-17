package router

import (
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
		query := p.Req.URL.Query().Get("q")
		fail.Render(p, pages.BookMgmtSearchFull(query))
	case "search-grid":
		query := p.Req.URL.Query().Get("q")
		fail.Render(p, pages.BookMgmtSearchGrid(query))

	default:
		fail.Redirect(p)
	}
}
