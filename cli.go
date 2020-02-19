package swap

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charliekenney23/tf-provider-swap/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// subcommands
var (
	cmdPreset = &cli.Command{
		Name:  "preset",
		Usage: "Manage presets",
		Subcommands: []*cli.Command{
			cmdPresetAdd,
			cmdPresetExec,
			cmdPresetList,
			cmdPresetRemove,
		},
	}

	cmdPresetAdd = &cli.Command{
		Name:  "add",
		Usage: "Add a preset",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "name of the preset",
				Aliases:  []string{"n"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "provider",
				Usage:    "provider to update",
				Aliases:  []string{"p"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "binPath",
				Usage:    "path to bin",
				Aliases:  []string{"bin", "b"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "preUpdate",
				Usage:   "shell command to run before swapping",
				Aliases: []string{"pre"},
			},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")
			binPath := c.String("binPath")

			absBinPath, err := filepath.Abs(binPath)
			if err != nil {
				return fmt.Errorf(`Error resolving binary file "%s": %s`, absBinPath, err.Error())
			}

			p := config.Preset{
				Provider:  c.String("provider"),
				PreUpdate: c.String("preUpdate"),
				BinPath:   absBinPath,
			}

			conf := config.Provider.GetConfig()
			if _, ok := conf.Presets[name]; ok {
				return fmt.Errorf(`Preset "%s" already exists`, name)
			}

			conf.Presets[name] = p
			config.Provider.SetConfig(conf)
			return nil
		},
	}

	cmdPresetRemove = &cli.Command{
		Name:      "remove",
		Usage:     "Remove a preset",
		ArgsUsage: "[name]",
		Action: func(c *cli.Context) error {
			name := c.Args().First()
			if name == "" {
				return errors.New("Name is required")
			}

			conf := config.Provider.GetConfig()
			if _, ok := conf.Presets[name]; !ok {
				return fmt.Errorf(`Preset "%s" does not exist`, name)
			}

			delete(conf.Presets, name)
			config.Provider.SetConfig(conf)
			return nil
		},
	}

	cmdPresetExec = &cli.Command{
		Name:  "exec",
		Usage: "Execute a preset",
		Action: func(c *cli.Context) error {
			AssertInTerraformWorkspace()

			name := c.Args().First()
			if name == "" {
				return errors.New("Preset is required")
			}

			conf := config.Provider.GetConfig()
			p, ok := conf.Presets[name]
			if !ok {
				return fmt.Errorf(`Preset "%s" does not exist`, name)
			}

			if p.PreUpdate != "" {
				color.Green("Running pre update script...")
				if err := runCommand(p.PreUpdate); err != nil {
					return fmt.Errorf("Error running pre-update script: %s", err.Error())
				}
			}

			if p.Provider == "" {
				return errors.New("Provider is required")
			}
			if p.BinPath == "" {
				return errors.New("Bin is required")
			}
			return UpdateProvider(p.Provider, p.BinPath)
		},
	}

	cmdPresetList = &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List all presets",
		Action: func(c *cli.Context) error {
			conf := config.Provider.GetConfig()
			for p := range conf.Presets {
				fmt.Println(p)
			}

			return nil
		},
	}

	cmdSwap = &cli.Command{
		Name:      "swap",
		Aliases:   []string{"s"},
		ArgsUsage: "[provider] [bin]",
		Usage:     "Swaps terraform provider binary for local builds",
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
)

// Entrypoint is the entrypoint of the CLI.
func Entrypoint() error {
	app := &cli.App{
		Name:                  "tf-provider-swap",
		Usage:                 "swap Terraform provider binaries in workspaces for local development.",
		Version:               "0.0.0",
		Compiled:              time.Now(),
		CustomAppHelpTemplate: helpTemplate,
		Commands: []*cli.Command{
			cmdSwap,
			cmdPreset,
		},
	}
	return app.Run(os.Args)
}
