package plugin

import (
	"fmt"

	wp "github.com/thegeeklab/wp-plugin-go/v2/plugin"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"
)

//go:generate go run ../internal/docs/main.go -output=../docs/data/data-raw.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*wp.Plugin
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	Action cli.StringSlice

	TofuVersion string
	InitOptions InitOptions
	FmtOptions  FmtOptions

	RootDir     string
	DataDir     string
	OutFile     string
	Parallelism int
	Targets     cli.StringSlice
	Refresh     bool
	NoLog       bool
}

// InitOptions include options for the OpenTofu init command.
type InitOptions struct {
	BackendConfig []string `json:"backend-config"`
	Lock          *bool    `json:"lock"`
	LockTimeout   string   `json:"lock-timeout"`
}

// FmtOptions fmt options for the OpenTofu fmt command.
type FmtOptions struct {
	List  *bool `json:"list"`
	Write *bool `json:"write"`
	Diff  *bool `json:"diff"`
	Check *bool `json:"check"`
}

type Cmd struct {
	*execabs.Cmd
	Private bool
}

func New(e wp.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := wp.Options{
		Name:                "wp-opentofu",
		Description:         "Manage infrastructure with OpenTofu",
		Flags:               Flags(p.Settings, wp.FlagsPluginCategory),
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

	p.Plugin = wp.New(options)

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
			Destination: &settings.NoLog,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "targets",
			Usage:       "targets to run `plan` or `apply` action on",
			EnvVars:     []string{"PLUGIN_TARGETS"},
			Destination: &settings.Targets,
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
			Destination: &settings.Refresh,
			Value:       true,
			Category:    category,
		},
	}
}
