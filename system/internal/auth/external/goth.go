package external

import (
	"fmt"
	"log"
	"sort"
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
	"github.com/crusttech/crust/internal/settings"
)

func setupGoth(c *config.Social, settings settings.Finder) {
	s, err := settings.FindByPrefix("auth.providers.")
	if err != nil {
		log.Printf("failed load 'auth.providers.' settings: %v", err)
	}

	if c == nil {
		log.Println("Skipping external auth setup")
		return
	}

	log.Println("Initializing external auth providers")

	// var defScopes = []string{"email"}

	store := sessions.NewCookieStore([]byte(c.SessionStoreSecret))
	store.MaxAge(c.SessionStoreExpiry)
	store.Options.Path = "/external"
	store.Options.HttpOnly = true
	store.Options.Secure = false // @todo, make this dependable on config somehow
	gothic.Store = store

	setupGothProviders(c.Url+"/external/%s/callback", s.KV())

	for p := range goth.GetProviders() {
		log.Printf("Social login initialized with %s provider", p)
	}
}

func setupGothProviders(callbackUrl string, s settings.KV) {
	// Purge all
	goth.ClearProviders()

	// Each setup...() func knows what to do
	setupGothFacebook(callbackUrl, s)
	setupGothGplus(callbackUrl, s)
	setupGothGithub(callbackUrl, s)
	setupGothLinkedin(callbackUrl, s)
	setupGothOpenIDConnect(callbackUrl, s)
}

func setupGothFacebook(callbackUrl string, s settings.KV) {
	const base = "auth.providers.facebook"
	if s.Bool(base) {
		goth.UseProviders(facebook.New(
			s.String(base+".key"),
			s.String(base+".secret"),
			fmt.Sprintf(callbackUrl, "facebook"),
			"email"))
	}
}

func setupGothGplus(callbackUrl string, s settings.KV) {
	const base = "auth.providers.gplus"
	if s.Bool(base) {
		goth.UseProviders(gplus.New(
			s.String(base+".key"),
			s.String(base+".secret"),
			fmt.Sprintf(callbackUrl, "gplus"),
			"email"))
	}
}

func setupGothGithub(callbackUrl string, s settings.KV) {
	const base = "auth.providers.github"
	if s.Bool(base) {
		goth.UseProviders(github.New(
			s.String(base+".key"),
			s.String(base+".secret"),
			fmt.Sprintf(callbackUrl, "github"),
			"email"))
	}
}

func setupGothLinkedin(callbackUrl string, s settings.KV) {
	const base = "auth.providers.linkedin"
	if s.Bool(base) {
		goth.UseProviders(linkedin.New(
			s.String(base+".key"),
			s.String(base+".secret"),
			fmt.Sprintf(callbackUrl, "linkedin"),
			"email"))
	}
}

// setupGothOpenIDConnect handles extracted OIDC providers
func setupGothOpenIDConnect(callbackUrl string, s settings.KV) {
	extractOpenIDProviders(s, func(name, url, clientKey, secret string) {
		if provider, err := openidConnect.New(clientKey, secret, callbackUrl, url, "email"); err != nil {
			log.Printf("failed to discover (auto discovery URL: %s) OIDC provider: %v", url, err)
		} else {
			provider.SetName(name)
			goth.UseProviders(provider)
		}
	})
}

// extractOpenIDProviders extracts OIDC providers and passes gathered settings to handler func
//
// Might
func extractOpenIDProviders(s settings.KV, handler func(name, url, clientKey, secret string)) {
	const base = "auth.providers.openid-connect"

	var (
		unique = map[string]bool{}
		name   string
	)

	for k := range s.Filter(base) {
		if len(k) < len(base)+2 {
			// skip invalid keys
			continue
		}

		// find next dot:
		name = k[len(base)+1:]
		dotPos := strings.Index(name, ".")
		if dotPos > 0 {
			name = name[:dotPos]
		}

		unique[name] = true
	}

	var sorted = make([]string, 0, len(unique))
	for gkey := range unique {
		sorted = append(sorted, gkey)
	}

	sort.Strings(sorted)

	for _, name = range sorted {
		// Enabled?
		if s.Bool(base + "." + name) {
			// Pass name, auto-discovery url, client key & secret to handler func
			handler(
				"openid-connect-"+name,
				s.String(base+"."+name+".url"),
				s.String(base+"."+name+".key"),
				s.String(base+"."+name+".secret"),
			)
		}
	}
}
