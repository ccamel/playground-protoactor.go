//nolint:gochecknoglobals,gochecknoinits // common pattern when using cobra library
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ccamel/playground-protoactor.go/internal/persistence/registry"
	"github.com/ccamel/playground-protoactor.go/internal/system"
)

var persistenceURI string

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the protoactor platform",
	Long:  `Start the protoactor platform`,
	RunE: func(_ *cobra.Command, _ []string) error {
		config, err := newSystemConfig()
		if err != nil {
			return err
		}

		sys, err := system.Boot(*config)
		if err != nil {
			return err
		}

		sys.Wait()

		return nil
	},
}

func init() {
	serveCmd.Flags().StringVar(
		&persistenceURI,
		"persistence-uri",
		"db:bbolt:./my-db?snapshotInterval=3",
		fmt.Sprintf("Persistence URI. Supported databases: %s.", strings.Join(registry.SupportedSchemes(), ", ")))

	rootCmd.AddCommand(serveCmd)
}

func newSystemConfig() (*system.Config, error) {
	var configOptions []system.Option

	if persistenceURI != "" {
		configOptions = append(configOptions, system.WithPersistenceURI(persistenceURI))
	}

	return system.NewConfig(configOptions...)
}
