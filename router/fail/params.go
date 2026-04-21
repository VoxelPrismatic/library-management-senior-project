package fail

import (
	"net/http"
	"strings"
	"voxelprismatic/library-management-senior-project/db"
)

type RoutingParams struct {
	W        http.ResponseWriter
	Req      *http.Request
	SubPtr   int
	FullPath []string
	User     *db.User
}

func (p *RoutingParams) Pop() string {
	if p.SubPtr >= len(p.FullPath) {
		return ""
	}

	ret := p.FullPath[p.SubPtr]
	p.SubPtr++
	return ret
}

func MakeParams(writer http.ResponseWriter, request *http.Request) *RoutingParams {
	u := db.CookieAuth(writer, request)
	p := strings.Split(request.URL.Path[len("/"):], "/")
	ret := RoutingParams{
		W:        writer,
		Req:      request,
		User:     &u,
		SubPtr:   0,
		FullPath: p,
	}
	return &ret
}

// Returns a joined path string from the current sub path pointer
func (p *RoutingParams) SubPath() string {
	return p.SubPathN(p.SubPtr)
}

// Returns a joined path string from a custom index.
// Use '0' for the full path.
func (p *RoutingParams) SubPathN(idx int) string {
	return "/" + strings.Join(p.FullPath[idx:], "/")
}

// Returns a joined path string until the current sub path pointer
// If 'include' is true, it'll include the current path part
func (p *RoutingParams) SubPathTree(include bool) string {
	plus := min(p.SubPtr, len(p.FullPath)-1)
	if include {
		plus++
	}
	return "/" + strings.Join(p.FullPath[:plus], "/")
}
