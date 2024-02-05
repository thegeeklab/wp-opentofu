package main

import (
	"github.com/thegeeklab/wp-opentofu/plugin"
	"github.com/urfave/cli/v2"
)

// SettingsFlags has the cli.Flags for the plugin.Settings.
//
//go:generate go run docs.go flags.go
func settingsFlags(settings *plugin.Settings, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "action",
			Usage:       "tofu actions to execute",
			EnvVars:     []string{"PLUGIN_ACTION"},
			Value:       cli.NewStringSlice("validate", "plan", "apply"),
			Destination: &settings.Action,
			Category:    category,
		},
		&cli.StringFlag{
			Name:     "init-options",
			Usage:    "tofu init command options, see https://opentofu.org/docs/cli/commands/init/",
			EnvVars:  []string{"PLUGIN_INIT_OPTIONS"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "fmt-options",
			Usage:    "options for the fmt command, see https://opentofu.org/docs/cli/commands/fmt/",
			EnvVars:  []string{"PLUGIN_FMT_OPTIONS"},
			Category: category,
		},
		&cli.IntFlag{
			Name:     "parallelism",
			Usage:    "number of concurrent operations",
			EnvVars:  []string{"PLUGIN_PARALLELISM"},
			Category: category,
		},
		&cli.StringFlag{
			Name:        "root-dir",
			Usage:       "root directory where the tofu files live",
			EnvVars:     []string{"PLUGIN_ROOT_DIR"},
			Destination: &settings.RootDir,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "no-log",
			Usage:       "suppress tofu command output",
			EnvVars:     []string{"PLUGIN_NO_LOG"},
			Destination: &settings.NoLog,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "targets",
			Usage:       "targets to run `apply` or `plan` action on",
			EnvVars:     []string{"PLUGIN_TARGETS"},
			Destination: &settings.Targets,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "version",
			Usage:       "tofu version to use",
			EnvVars:     []string{"PLUGIN_VERSION"},
			Destination: &settings.Version,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "refresh",
			Usage:       "enables refreshing of the state before `plan` and `apply` commands",
			EnvVars:     []string{"PLUGIN_REFRESH"},
			Destination: &settings.Refresh,
			Value:       true,
			Category:    category,
		},
	}
}
