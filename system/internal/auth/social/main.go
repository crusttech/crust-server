package social

import (
	"github.com/crusttech/crust/internal/config"
)

func Init(c *config.Social) {
	initGoth(c)
}
