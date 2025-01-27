package handler

import (
	"net/http"

	"github.com/caos/zitadel/internal/domain"
)

const (
	tmplLoginSuccess = "login_success"
)

type loginSuccessData struct {
	userData
	RedirectURI string `schema:"redirect-uri"`
}

func (l *Login) redirectToLoginSuccess(w http.ResponseWriter, r *http.Request, id string) {
	http.Redirect(w, r, l.renderer.pathPrefix+EndpointLoginSuccess+"?authRequestID="+id, http.StatusFound)
}

func (l *Login) handleLoginSuccess(w http.ResponseWriter, r *http.Request) {
	authRequest, _ := l.getAuthRequest(r)
	if authRequest == nil {
		l.renderSuccessAndCallback(w, r, nil, nil)
		return
	}
	for _, step := range authRequest.PossibleSteps {
		if step.Type() != domain.NextStepLoginSucceeded && step.Type() != domain.NextStepRedirectToCallback {
			l.renderNextStep(w, r, authRequest)
			return
		}
	}
	l.renderSuccessAndCallback(w, r, authRequest, nil)
}

func (l *Login) renderSuccessAndCallback(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := loginSuccessData{
		userData: l.getUserData(r, authReq, "Login Successful", errID, errMessage),
	}
	if authReq != nil {
		data.RedirectURI = l.oidcAuthCallbackURL
	}
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplLoginSuccess], data, nil)
}

func (l *Login) redirectToCallback(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	callback := l.oidcAuthCallbackURL + authReq.ID
	http.Redirect(w, r, callback, http.StatusFound)
}
