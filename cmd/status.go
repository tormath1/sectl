package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/tormath1/sectl/pkg"
)

var (
	// configFile is the file holding the
	// SELinux config
	configFile string

	statusCmd = &cobra.Command{
		Use:     "status",
		Short:   "display the status of SELinux",
		Example: "sectl status -o json",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := afero.NewOsFs()

			c, err := pkg.GetStatus(fs, configFile)
			if err != nil {
				return fmt.Errorf("unable to read config: %w", err)
			}

			output := cmd.Flag("output")
			// TODO: export this into a dedicated `print` package
			switch output.Value.String() {
			case "json":
				res, err := json.Marshal(c)
				if err != nil {
					return fmt.Errorf("unable to generate JSON: %w", err)
				}

				fmt.Fprint(os.Stdout, string(res))
			}

			return nil
		},
	}
)

func init() {
	statusCmd.Flags().StringVarP(&configFile, "config-file", "c", "/etc/selinux/config", "path of the config file")
}
