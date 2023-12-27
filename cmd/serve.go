//nolint:gochecknoglobals,gochecknoinits // common pattern when using cobra library
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ccamel/playground-protoactor.go/internal/system"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the protoactor platform",
	Long:  `Start the protoactor platform`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sys, err := system.Boot()
		if err != nil {
			return err
		}

		sys.Wait()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
