package swap

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"
)

// Entrypoint is the entrypoint of the CLI.
func Entrypoint() error {
	app := &cli.App{
		Name:      "tf-provider-swap",
		Usage:     "Swap Terraform providers for local builds",
		ArgsUsage: "[provider] [bin]",
		Action: func(c *cli.Context) error {
			AssertInTerraformWorkspace()

			argv := c.Args()
			provider := argv.Get(0)
			binPath := argv.Get(1)

			if provider == "" {
				return errors.New("Provider is required")
			}
			if binPath == "" {
				return errors.New("Bin is required")
			}
			return UpdateProvider(provider, binPath)
		},
	}

	return app.Run(os.Args)
}
