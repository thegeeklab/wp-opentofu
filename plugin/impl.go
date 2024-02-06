package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/execabs"
)

var (
	ErrTaintedPath        = errors.New("filepath is tainted")
	ErrMaxSizeSizeLimit   = errors.New("max size limit of decoded data exceeded")
	ErrActionUnknown      = errors.New("action not found")
	ErrInvalidTofuVersion = errors.New("invalid version string")
)

const (
	maxDecompressionSize = 1024
	defaultDirPerm       = 0o755
)

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	p.Settings.DataDir = ".terraform"
	if value, ok := os.LookupEnv("TF_DATA_DIR"); ok {
		p.Settings.DataDir = value
	}

	p.Settings.OutFile = "plan.tfout"
	if p.Settings.DataDir == ".terraform" {
		p.Settings.OutFile = fmt.Sprintf("%s.plan.tfout", p.Settings.DataDir)
	}

	if p.Settings.TofuVersion != "" {
		err := installPackage(p.Plugin.Network.Context, p.Plugin.Network.Client, p.Settings.TofuVersion, maxDecompressionSize)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	commands := []*execabs.Cmd{
		p.versionCommand(),
	}

	commands = append(commands, p.initCommand())
	commands = append(commands, p.getModulesCommand())

	for _, action := range p.Settings.Action.Value() {
		switch action {
		case "fmt":
			commands = append(commands, p.fmtCommand())
		case "validate":
			commands = append(commands, p.validateCommand())
		case "plan":
			commands = append(commands, p.planCommand(false))
		case "plan-destroy":
			commands = append(commands, p.planCommand(true))
		case "apply":
			commands = append(commands, p.applyCommand())
		case "destroy":
			commands = append(commands, p.destroyCommand())
		default:
			return fmt.Errorf("%w: %s", ErrActionUnknown, action)
		}
	}

	if err := deleteCache(p.Settings.DataDir); err != nil {
		return err
	}

	for _, cmd := range commands {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if p.Settings.RootDir != "" {
			cmd.Dir = p.Settings.RootDir
		}

		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return deleteCache(p.Settings.DataDir)
}
