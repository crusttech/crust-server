package social

import (
	"log"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/openidConnect"

	"github.com/crusttech/crust/internal/config"
)

func initGoth(c *config.Social) {
	if c == nil {
		log.Println("Skipping social login setup")
		return
	}

	log.Println("Initializing social login providers")

	var defScopes = []string{"email"}

	store := sessions.NewCookieStore([]byte(c.SessionStoreSecret))
	store.MaxAge(c.SessionStoreExpiry)
	store.Options.Path = "/social"
	store.Options.HttpOnly = true
	store.Options.Secure = false // @todo
	gothic.Store = store

	// Returns false if any of the passed string values is empty
	has := func(ss ...string) bool {
		for _, s := range ss {
			if strings.TrimSpace(s) == "" {
				return false
			}
		}

		return true
	}

	if has(c.OidcUrl) {
		if provider, err := openidConnect.New(c.OidcKey, c.OidcSecret, c.Url+"/social/openid-connect/callback", c.OidcUrl, defScopes...); err != nil {
			log.Printf("failed to discover (auto discovery URL: %s) OIDC provider: %v", c.OidcUrl, err)
		} else {
			goth.UseProviders(provider)
		}
	}

	if has(c.FacebookKey, c.FacebookSecret) {
		goth.UseProviders(facebook.New(c.FacebookKey, c.FacebookSecret, c.Url+"/social/facebook/callback", defScopes...))
	}

	if has(c.GPlusKey, c.GPlusSecret) {
		goth.UseProviders(gplus.New(c.GPlusKey, c.GPlusSecret, c.Url+"/social/gplus/callback", defScopes...))
	}

	if has(c.GitHubKey, c.GitHubSecret) {
		goth.UseProviders(github.New(c.GitHubKey, c.GitHubSecret, c.Url+"/social/github/callback", defScopes...))
	}

	if has(c.LinkedInKey, c.LinkedInSecret) {
		goth.UseProviders(linkedin.New(c.LinkedInKey, c.LinkedInSecret, c.Url+"/social/linkedin/callback", defScopes...))
	}

	for p := range goth.GetProviders() {
		log.Printf("Social login initialized with %s provider", p)
	}
}
