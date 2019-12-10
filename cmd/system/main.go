package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/crusttech/crust-server/pkg/subscription"
)

func main() {
	cfg := system.Configure()
	cfg.RootCommandName = "crust-server-system"

	cfg.ApiServerPreRun = append(
		cfg.ApiServerPreRun,
		func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
			if service.CurrentSubscription != nil {
				// Already initialized
				return nil
			}

			subscription.Init(logger.Default(), service.DefaultSettings)
			subscription.UpdateCurrent(subscription.Load(ctx))
			return nil
		},
	)

	cmd := cfg.MakeCLI(cli.Context())
	cli.HandleError(cmd.Execute())
}
