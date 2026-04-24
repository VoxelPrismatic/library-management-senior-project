package router

import (
	"fmt"
	"net/http"
	"voxelprismatic/library-management-senior-project/db"
	"voxelprismatic/library-management-senior-project/router/fail"
	"voxelprismatic/library-management-senior-project/web/user"
)

func UserRouter(p *fail.RoutingParams) {
	switch p.Pop() {
	case "register":
		HandleUserRegister(p)
	default:
		fail.Redirect(p)
	}
}

func HandleUserRegister(p *fail.RoutingParams) {
	if fail.Done(p) {
		return
	}

	switch p.Req.Method {
	case http.MethodGet:
		fail.Render(p, user.Register())
	case http.MethodPost:
		HandleUserRegisterPost(p)
	default:
		http.Error(p.W, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleUserRegisterPost(p *fail.RoutingParams) {
	if fail.Form(p) {
		return
	}

	userObj := db.User{}
	errs := map[string]error{}
	formData := p.Form()

	if err := userObj.SetFirstName(formData["firstName"]); err != nil {
		errs["firstName"] = err
	}

	if err := userObj.SetLastName(formData["lastName"]); err != nil {
		errs["lastName"] = err
	}

	if err := userObj.SetEmail(formData["emailAddr"]); err != nil {
		errs["emailAddr"] = err
	}

	if err := db.TestSecretStrength(formData["secret"]); err != nil {
		errs["secret"] = err
	}

	if err := userObj.SetSecret(formData["secret"], formData["secret_again"]); err != nil {
		errs["secret_again"] = err
	}

	if len(errs) > 0 {
		(fail.HTMX{
			Retarget: "#formEntry",
			Reswap:   "outerHTML",
		}).Apply(p)
		fail.Render(p, user.FormTable(formData, errs))
		return
	}

	jwt := userObj.IssueJWT()
	p.W.Header().Set("Set-Cookie", fmt.Sprintf("tok=%s", jwt.Token))
	(fail.HTMX{
		Redirect: "/",
	}).Apply(p)
}
