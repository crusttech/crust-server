package subscription

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/settings"
)

type (
	// Interface to settings backend
	//
	// We are using settings backend for storing subscription key
	//  - crust-subscription.jet
	settingsGetterSetter interface {
		Get(context.Context, string, uint64) (*settings.Value, error)
		Set(context.Context, *settings.Value) error
	}
)

const (
	settingSubscriptionJwtKey   = "crust-subscription.jwt"
	settingSubscriptionTrialKey = "crust-subscription.trial"
)

var (
	logger = zap.NewNop()

	publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAqt/kkmAzRJ+OMC2aIN61En2syeiu
GboKvT/xEbbmEXXopvS/Pcwrgd/Ml+O/+AxjNcfF9c5De7FAs9TdkTo26JQB/Y3a
IhMeAVbGZ1dihYuMKeubc/StTDdwWb4UvhT6cmZl559zpU7f06LWTSvAtM+MOK+A
SzbZ8RTdJSvY75iw0TI=
-----END PUBLIC KEY-----`)

	settingsSvc settingsGetterSetter
)

// Init sets pkg basics: logger & settings interface
func Init(l *zap.Logger, ss settingsGetterSetter) {
	logger = l.Named("crust-subscription").
		// Do not attach stack trace to any of logs below DPanic
		//
		// None of the logs at subscription pkg should be at DPanic or higher.
		WithOptions(zap.AddStacktrace(zap.DPanicLevel))

	if ss == nil {
		logger.Error("could not load subscription settings, no settings service")
		return
	}

	settingsSvc = ss
}

func Load(ctx context.Context) *Claims {
	if v, err := settingsSvc.Get(ctx, settingSubscriptionJwtKey, 0); err != nil {
		logger.Error("could not load subscription JWT key", zap.Error(err))
		return nil
	} else if v == nil || v.String() == "" {
		logger.Info("subscription value missing", zap.String("name", settingSubscriptionJwtKey))
		return genericTrial(ctx)
	} else {
		return parse(v.String())
	}
}

// Parses subscription
func parse(subval string) *Claims {
	var claims = &Claims{}

	parsedToken, err := jwt.ParseWithClaims(subval, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwt.ParseECPublicKeyFromPEM(publicKey)
	})

	if err != nil {
		logger.Error("failed to parse subscription jwt", zap.Error(err))
		return nil
	}

	if !parsedToken.Valid || parsedToken.Header["type"] != HEADER_TYPE {
		logger.Error("invalid subscription jwt")
		return nil
	}

	if err = claims.Valid(); err != nil {
		logger.Error("invalid subscription", zap.Error(err))
		return nil
	}

	logger.Debug("subscription loaded")

	return claims
}

// Generates trial if it does not exist yet
//
// this is a simple date entry in the setting table -- YYYY-MM-DD when
// trial expires
//
// This function is called only when no (other) subscription key setting is found
func genericTrial(ctx context.Context) *Claims {
	const expDateFormat = "2006-01-02"

	var (
		saveTrial bool
		expDate   time.Time
	)

	if v, err := settingsSvc.Get(ctx, settingSubscriptionTrialKey, 0); err != nil {
		logger.Error("could not load subscription trial", zap.Error(err))
		return nil
	} else if v != nil || v.String() != "" {
		logger.Info("generic trial key found")

		// Trial exp. date found
		expDate, err = time.Parse(expDateFormat, v.String())
		if err != nil {
			// If we can not parse the trial exp. date, assume it expired today.
			expDate = time.Date(now().Year(), now().Month(), now().Day(), 0, 0, 0, 0, now().Location())
			saveTrial = true
		}

	} else {
		logger.Info("creating generic trial key")

		// Setting not found, create 30 day expiration
		expDate = time.Date(now().Year(), now().Month(), now().Day()+31, 0, 0, 0, 0, now().Location())
		saveTrial = true
	}

	if saveTrial {
		v := &settings.Value{Name: settingSubscriptionTrialKey}
		_ = v.SetValue(expDate.Format(expDateFormat))
		if err := settingsSvc.Set(ctx, v); err != nil {
			logger.Error("could not save subscription trial", zap.Error(err))
			return nil
		}
	}

	return &Claims{Trial: true, MaxUsers: 10, Expires: expDate}
}
