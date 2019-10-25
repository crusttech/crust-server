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

	// initServices := cfg.InitServices
	// cfg.InitServices = func(ctx context.Context, c *cli.Config) {
	// 	// subscription.Check(context.Background(), "local.cortezaproject.org", "")
	// 	initServices(ctx, c)
	// 	subscription.Init(logger.Default(), service.DefaultIntSettings)
	// }

	cfg.ApiServerPreRun = append(
		cfg.ApiServerPreRun,
		func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
			// subscription.Check(context.Background(), "local.cortezaproject.org", "")
			subscription.Init(ctx, logger.Default(), service.DefaultIntSettings)
			return nil
		},
	)

	cmd := cfg.MakeCLI(cli.Context())
	cli.HandleError(cmd.Execute())
}
