package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/thegeeklab/wp-plugin-go/v2/trace"
)

var (
	ErrTaintedPath        = errors.New("filepath is tainted")
	ErrMaxSizeSizeLimit   = errors.New("max size limit of decoded data exceeded")
	ErrActionUnknown      = errors.New("action not found")
	ErrInvalidTofuVersion = errors.New("invalid version string")
	ErrHTTPError          = errors.New("http error")
)

const (
	maxDecompressionSize = 100 * 1024 * 1024
	defaultDirPerm       = 0o755
)

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
	if err := p.FlagsFromContext(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (p *Plugin) FlagsFromContext() error {
	if p.Context.String("init-option") != "" {
		initOptions := InitOptions{}
		if err := json.Unmarshal([]byte(p.Context.String("init-option")), &initOptions); err != nil {
			return fmt.Errorf("cannot unmarshal init_option: %w", err)
		}

		p.Settings.InitOptions = initOptions
	}

	if p.Context.String("fmt-option") != "" {
		fmtOptions := FmtOptions{}
		if err := json.Unmarshal([]byte(p.Context.String("fmt-option")), &fmtOptions); err != nil {
			return fmt.Errorf("cannot unmarshal fmt_option: %w", err)
		}

		p.Settings.FmtOptions = fmtOptions
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

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	batchCmd := make([]*Cmd, 0)
	batchCmd = append(batchCmd, p.versionCommand())

	if p.Settings.TofuVersion != "" {
		err := installPackage(p.Plugin.Network.Context, p.Plugin.Network.Client, p.Settings.TofuVersion, maxDecompressionSize)
		if err != nil {
			return err
		}
	}

	batchCmd = append(batchCmd, p.initCommand())
	batchCmd = append(batchCmd, p.getModulesCommand())

	for _, action := range p.Settings.Action.Value() {
		switch action {
		case "fmt":
			batchCmd = append(batchCmd, p.fmtCommand())
		case "validate":
			batchCmd = append(batchCmd, p.validateCommand())
		case "plan":
			batchCmd = append(batchCmd, p.planCommand(false))
		case "plan-destroy":
			batchCmd = append(batchCmd, p.planCommand(true))
		case "apply":
			batchCmd = append(batchCmd, p.applyCommand())
		case "destroy":
			batchCmd = append(batchCmd, p.destroyCommand())
		default:
			return fmt.Errorf("%w: %s", ErrActionUnknown, action)
		}
	}

	if err := deleteDir(p.Settings.DataDir); err != nil {
		return err
	}

	for _, bc := range batchCmd {
		bc.Stdout = os.Stdout
		bc.Stderr = os.Stderr
		trace.Cmd(bc.Cmd)

		bc.Env = os.Environ()

		if bc.Private {
			bc.Stdout = io.Discard
		}

		if p.Settings.RootDir != "" {
			bc.Dir = p.Settings.RootDir
		}

		if err := bc.Run(); err != nil {
			return err
		}
	}

	return deleteDir(p.Settings.DataDir)
}
