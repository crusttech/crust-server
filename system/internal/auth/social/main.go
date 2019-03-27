package social

import (
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/settings"
)

func Init(c *config.Social, finder settings.Finder) {
	setupGoth(c, finder)
}
