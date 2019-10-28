package subscription

import (
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/system/service"
)

type (
	Claims struct {
		Domains  []string
		Trial    bool
		MaxUsers uint
		Expires  time.Time
	}
)

const (
	HEADER_TYPE = "crust-subscription"
)

func (Claims) Valid() error {
	return nil
}

func UpdateCurrent(c *Claims) {
	if c == nil {
		return
	}

	if service.CurrentSubscription == nil {
		service.CurrentSubscription = &subscription{}
	}

	if s, ok := service.CurrentSubscription.(*subscription); !ok {
		logger.Error("unknown service.CurrentSubscription type")
	} else {
		s.Update(c.Domains, c.Expires, c.Trial, c.MaxUsers)

		logger.Info("subscription updated",
			zap.Strings("domains", c.Domains),
			zap.Time("expires", c.Expires),
			zap.Bool("is-trial", c.Trial),
			zap.Uint("limit-max-users", c.MaxUsers))
	}
}
