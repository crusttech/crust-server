package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/system/internal/service"
)

type (
	Social struct {
		auth       service.AuthService
		config     *config.Social
		jwtEncoder auth.TokenEncoder
	}
)

func NewSocial(config *config.Social, jwtEncoder auth.TokenEncoder) *Social {
	return &Social{
		auth:       service.DefaultAuth,
		config:     config,
		jwtEncoder: jwtEncoder,
	}
}

func (ctrl *Social) MountRoutes(r chi.Router) {
	// Copy provider from path (Chi URL param) to request context and return it
	copyProviderToContext := func(r *http.Request) *http.Request {
		return r.WithContext(context.WithValue(r.Context(), "provider", chi.URLParam(r, "provider")))
	}

	r.Route("/social/{provider}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			r = copyProviderToContext(r)

			// Always set redir cookie, even if not requested. If param is empty, cookie is removed
			ctrl.setCookie(w, r, "redir", r.URL.Query().Get("redir"))

			// try to get the user without re-authenticating
			if user, err := gothic.CompleteUserAuth(w, r); err != nil {
				gothic.BeginAuthHandler(w, r)
			} else {
				// We've successfully singed-in through 3rd party auth
				ctrl.handleSuccessfulAuth(w, r, user)
			}
		})

		r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
			r = copyProviderToContext(r)

			if user, err := gothic.CompleteUserAuth(w, r); err != nil {
				log.Printf("Failed to complete user auth: %v", err)
				ctrl.handleFailedCallback(w, r, err)
			} else {
				ctrl.handleSuccessfulAuth(w, r, user)
			}
		})

		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			if err := gothic.Logout(w, r); err != nil {
				log.Printf("Failed to social logout: %v", err)
			}

			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
		})
	})
}

func (ctrl *Social) handleFailedCallback(w http.ResponseWriter, r *http.Request, err error) {
	provider := chi.URLParam(r, "provider")

	if strings.Contains(err.Error(), "Error processing your OAuth request: Invalid oauth_verifier parameter") {
		// Just take user through the same loop again
		w.Header().Set("Location", "/social/"+provider)
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w, "SSO Error: %v", err.Error())
	w.WriteHeader(http.StatusOK)
}

// Handles authentication via external auth providers of
// unknown an user + appending authentication on external providers
// to a current user
func (ctrl *Social) handleSuccessfulAuth(w http.ResponseWriter, r *http.Request, cred goth.User) {
	log.Printf("Successful social login: %v", cred)

	if u, err := ctrl.auth.With(r.Context()).Social(cred); err != nil {
		resputil.JSON(w, err)
	} else {
		ctrl.jwtEncoder.SetCookie(w, r, u)

		if c, err := r.Cookie("redir"); c != nil && err == nil {
			ctrl.setCookie(w, r, "redir", "")
			w.Header().Set("Location", c.Value)
			w.WriteHeader(http.StatusSeeOther)

		}

		resputil.JSON(w, u, err)
	}
}

// Extracts and authenticates JWT from context, validates claims
func (ctrl *Social) setCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := &http.Cookie{
		Name:    name,
		Expires: time.Now().Add(time.Duration(ctrl.config.SessionStoreExpiry) * time.Second),
		Secure:  r.URL.Scheme == "https",
		Domain:  r.URL.Hostname(),
		Path:    "/social",
	}

	if value == "" {
		cookie.Expires = time.Unix(0, 0)
	} else {
		cookie.Value = value
	}

	http.SetCookie(w, cookie)
}
