package plugin

import (
	"fmt"

	"github.com/thegeeklab/wp-opentofu/tofu"
	plugin_base "github.com/thegeeklab/wp-plugin-go/v4/plugin"
	"github.com/urfave/cli/v2"
)

//go:generate go run ../internal/docs/main.go -output=../docs/data/data-raw.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*plugin_base.Plugin
	*plugin_base.Environment
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	Action      cli.StringSlice
	RootDir     string
	DataDir     string
	TofuVersion string
	Tofu        tofu.Tofu
}

func New(e plugin_base.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := plugin_base.Options{
		Name:                "wp-opentofu",
		Description:         "Manage infrastructure with OpenTofu",
		Flags:               Flags(p.Settings, plugin_base.FlagsPluginCategory),
		Execute:             p.run,
		HideWoodpeckerFlags: true,
	}

	if len(build) > 0 {
		options.Version = build[0]
	}

	if len(build) > 1 {
		options.VersionMetadata = fmt.Sprintf("date=%s", build[1])
	}

	if e != nil {
		options.Execute = e
	}

	p.Plugin = plugin_base.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
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
			Name:     "init-option",
			Usage:    "tofu init command options, see https://opentofu.org/docs/cli/commands/init/",
			EnvVars:  []string{"PLUGIN_INIT_OPTION"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "fmt-option",
			Usage:    "options for the fmt command, see https://opentofu.org/docs/cli/commands/fmt/",
			EnvVars:  []string{"PLUGIN_FMT_OPTION"},
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
			Usage:       "suppress tofu command output for `plan`, `apply` and `destroy` action",
			EnvVars:     []string{"PLUGIN_NO_LOG"},
			Destination: &settings.Tofu.NoLog,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "targets",
			Usage:       "targets to run `plan` or `apply` action on",
			EnvVars:     []string{"PLUGIN_TARGETS"},
			Destination: &settings.Tofu.Targets,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "tofu-version",
			Usage:       "tofu version to use",
			EnvVars:     []string{"PLUGIN_TOFU_VERSION"},
			Destination: &settings.TofuVersion,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "refresh",
			Usage:       "enables refreshing of the state before `plan` and `apply` commands",
			EnvVars:     []string{"PLUGIN_REFRESH"},
			Destination: &settings.Tofu.Refresh,
			Value:       true,
			Category:    category,
		},
	}
}
