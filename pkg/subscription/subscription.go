package subscription

import (
	"context"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/service"
)

type (
	// Interface to settings backend
	//
	// We are using settings backend for storing domain & key:
	//  - crust-subscription.domain
	//  - crust-subscription.key
	settingsGetter interface {
		Get(string, uint64) (*settings.Value, error)
		Set(*settings.Value) error
	}
)

var (
	settingsService settingsGetter
)

// Init sets up crust subscriptions, sets logger and settings getter.
func Init(ctx context.Context, l *zap.Logger, s settingsGetter) {
	logger = l.Named("crust-subscription").
		// Do not attach stack trace to any of logs below DPanic
		//
		// None of the logs at subscription pkg should be at DPanic or higher.
		WithOptions(zap.AddStacktrace(zap.DPanicLevel))
	settingsService = s

	if service.CurrentSubscription != nil {
		// Already initialized
		return
	}

	// make sure we do not run 2x
	service.CurrentSubscription = &subscription{}

	var domain = getValue("domain")
	if domain == "" {
		logger.Error("refusing to initialize subscription, missing domain")
		return
	}

	logger.Info("loading subscription for domain", zap.String("domain", domain))
	service.CurrentSubscription = New(domain)

	if key := getValue("key"); key != "" {
		// We have subscription key, update current subscription
		updateCurrent(ctx, domain, key)
	} else if email := getValue("email"); email != "" {
		// no subscription key, but we have email
		requestTrial(ctx, domain, email)
	} else {
		logger.Warn("could not check subscription or create trial: missing key and email values")
	}
}

func UpdateCurrent(ctx context.Context) {
	if service.CurrentSubscription != nil {
		// Not initialized, nothing to do here
		logger.Warn("current subscription not initialized")
		return
	}

	domain, key := getValue("domain"), getValue("key")
	if domain == "" {
		logger.Error("refusing to update subscription, missing domain")
		return
	} else if key == "" {
		logger.Error("refusing to update subscription, missing key")
		return
	}

	updateCurrent(ctx, domain, key)
}

func updateCurrent(ctx context.Context, domain, key string) {
	cs, ok := service.CurrentSubscription.(*subscription)
	if !ok {
		logger.Error("refusing to update foreign subscription")
		return
	}

	rsp, err := check(ctx, domain, key)
	if err != nil {
		logger.Error("could not check subscription", zap.Error(err))
	}

	cs.Update(rsp.Expires, rsp.IsTrial, rsp.LimitMaxUsers)

	logger.Info("subscription updated",
		zap.String("domain", cs.domain),
		zap.Time("expires", cs.expires),
		zap.Bool("is-trial", cs.isTrial),
		zap.Uint("limit-max-users", cs.limitMaxUsers))
}

func requestTrial(ctx context.Context, domain, email string) {
	cs, ok := service.CurrentSubscription.(*subscription)
	if !ok {
		logger.Error("refusing to request trial on a foreign subscription")
		return
	}

	rsp, err := trial(ctx, domain, email)
	if err != nil {
		logger.Error("could not request trial subscription", zap.Error(err))
		return
	}

	storeSetting("key", rsp.Key)
	storeSetting("domain", domain)

	cs.Update(rsp.Expires, rsp.IsTrial, rsp.LimitMaxUsers)

	logger.Info("trial subscription created",
		zap.String("domain", cs.domain),
		zap.String("email", email),
		zap.Time("expires", cs.expires),
		zap.Uint("limit-max-users", cs.limitMaxUsers))
}

// getValue tries to find subscription value (key, email, domain) in settings or in environmental variables
//
// It prefixes name with "crust-subscription." for settings and with "CRUST_SUBSCRIPTION_" for env-var lookup
func getValue(name string) (val string) {
	if settingsService == nil {
		logger.Error("could not load subscription settings, no settings service")
		return
	}

	var (
		settingKey = "crust-subscription." + name
		envKey     = "CRUST_SUBSCRIPTION_" + strings.ToUpper(name)
	)

	if v, err := settingsService.Get(settingKey, 0); err != nil {
		logger.Error("could not load subscription settings", zap.String("name", settingKey), zap.Error(err))
		return
	} else if v != nil || v.String() != "" {
		return v.String()
	} else if v, has := os.LookupEnv(envKey); has {
		return v
	} else if name == "domain" {
		if val = envHostnameLookup(); val != "" {
			return val
		}
	}

	logger.Info("subscription value missing",
		zap.String("settings-key", settingKey),
		zap.String("env-key", envKey),
	)

	return

}

func storeSetting(name, value string) bool {
	if settingsService == nil {
		logger.Error("could not store subscription settings, no settings service")
		return false
	}

	v := &settings.Value{Name: "crust-subscription." + name}
	if err := v.SetValue(value); err != nil {
		logger.Error("could not set subscription setting value", zap.String("name", name), zap.Error(err))
		return false
	}

	if err := settingsService.Set(v); err != nil {
		logger.Error("could not save subscription settings", zap.String("name", name), zap.Error(err))
		return false
	}

	return true
}

// Scans list of env variables, returns first one with valid value
func envHostnameLookup() string {
	// All env keys we'll check, first that has any value set, will be used as domain
	keysWithHostnames := []string{"DOMAIN", "LETSENCRYPT_HOST", "VIRTUAL_HOST", "HOSTNAME", "HOST"}

	for _, key := range keysWithHostnames {
		logger.Debug("searching for host in environmental variables",
			zap.String("key", key))

		if value, has := os.LookupEnv(key); has && len(value) > 0 {
			logger.Info("found host in environmental variable",
				zap.String("key", key),
				zap.String("value", value))
			return value
		}
	}

	return ""
}
