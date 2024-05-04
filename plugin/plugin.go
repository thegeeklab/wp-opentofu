package plugin

import (
	wp "github.com/thegeeklab/wp-plugin-go/v2/plugin"
	"github.com/urfave/cli/v2"
)

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

func New(options wp.Options, settings *Settings) *Plugin {
	p := &Plugin{}

	if options.Execute == nil {
		options.Execute = p.run
	}

	p.Plugin = wp.New(options)
	p.Settings = settings

	return p
}
