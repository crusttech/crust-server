package main

import (
	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/corteza"
	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/system"

	"github.com/cortezaproject/corteza-server/monolith"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func main() {
	logger.Init()

	app.Run(
		logger.Default(),
		app.NewOptions(),
		&corteza.App{},
		&monolith.App{
			System:    &system.App{},
			Compose:   &compose.App{},
			Messaging: &messaging.App{},
		},
	)
}
