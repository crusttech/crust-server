package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/settings"
)

func Settings(rootCmd *cobra.Command, service settings.Service) {
	exit := func(err error) {
		if err != nil {
			rootCmd.Printf("Error: %v\n", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	settingsCmd := &cobra.Command{
		Use:   "settings",
		Short: "Settings management",
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "List all",
		Run: func(cmd *cobra.Command, args []string) {
			prefix := cmd.Flags().Lookup("prefix").Value.String()
			if kv, err := service.FindByPrefix(prefix); err != nil {
				exit(err)
			} else {
				for k, v := range kv {
					cmd.Printf("%s\t%v\n", k, v)
				}
			}
		},
	}

	list.Flags().String("prefix", "", "Filter settings by prefix")

	get := &cobra.Command{
		Use: "get [key to get]",

		Short: "Get value (raw JSON) for a specific key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := service.Get(args[0], 0); err != nil {
				exit(err)
			} else if v != nil {
				cmd.Printf("%v\n", v.Value)
			}
			exit(nil)
		},
	}

	settingsCmd.AddCommand(list, get)

	settingsCmd.AddCommand(&cobra.Command{
		Use:   "set [key to set] [value if not stdin]",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			value := args[1]
			v := &settings.Value{
				Name: args[0],
			}

			if err := v.SetValueAsString(value); err != nil {
				exit(err)
			}

			exit(service.Set(v))
		},
	})

	rootCmd.AddCommand(settingsCmd)
}
