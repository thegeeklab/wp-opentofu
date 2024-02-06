package plugin

import (
	"fmt"

	"github.com/thegeeklab/wp-plugin-go/trace"
	"golang.org/x/sys/execabs"
)

const (
	tofuBin = "/usr/local/bin/tofu"
)

func (p *Plugin) versionCommand() *execabs.Cmd {
	args := []string{
		"version",
	}

	return execabs.Command(
		tofuBin,
		args...,
	)
}

func (p *Plugin) initCommand() *execabs.Cmd {
	args := []string{
		"init",
	}

	for _, v := range p.Settings.InitOptions.BackendConfig {
		args = append(args, fmt.Sprintf("-backend-config=%s", v))
	}

	// Fail tofu execution on prompt
	args = append(args, "-input=false")

	return execabs.Command(
		tofuBin,
		args...,
	)
}

func (p *Plugin) getModulesCommand() *execabs.Cmd {
	return execabs.Command(
		tofuBin,
		"get",
	)
}

func (p *Plugin) validateCommand() *execabs.Cmd {
	return execabs.Command(
		tofuBin,
		"validate",
	)
}

func (p *Plugin) fmtCommand() *execabs.Cmd {
	args := []string{
		"fmt",
	}

	if p.Settings.FmtOptions.List != nil {
		args = append(args, fmt.Sprintf("-list=%t", *p.Settings.FmtOptions.List))
	}

	if p.Settings.FmtOptions.Write != nil {
		args = append(args, fmt.Sprintf("-write=%t", *p.Settings.FmtOptions.Write))
	}

	if p.Settings.FmtOptions.Diff != nil {
		args = append(args, fmt.Sprintf("-diff=%t", *p.Settings.FmtOptions.Diff))
	}

	if p.Settings.FmtOptions.Check != nil {
		args = append(args, fmt.Sprintf("-check=%t", *p.Settings.FmtOptions.Check))
	}

	return execabs.Command(
		tofuBin,
		args...,
	)
}

func (p *Plugin) planCommand(destroy bool) *execabs.Cmd {
	args := []string{
		"plan",
	}

	if destroy {
		args = append(args, "-destroy")
	} else {
		args = append(args, fmt.Sprintf("-out=%s", p.Settings.OutFile))
	}

	for _, value := range p.Settings.Targets.Value() {
		args = append(args, "--target", value)
	}

	if p.Settings.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", p.Settings.Parallelism))
	}

	if p.Settings.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *p.Settings.InitOptions.Lock))
	}

	if p.Settings.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", p.Settings.InitOptions.LockTimeout))
	}

	if !p.Settings.Refresh {
		args = append(args, "-refresh=false")
	}

	cmd := execabs.Command(
		tofuBin,
		args...,
	)

	if !p.Settings.NoLog {
		trace.Cmd(cmd)
	}

	return cmd
}

func (p *Plugin) applyCommand() *execabs.Cmd {
	args := []string{
		"apply",
	}

	for _, v := range p.Settings.Targets.Value() {
		args = append(args, "--target", v)
	}

	if p.Settings.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", p.Settings.Parallelism))
	}

	if p.Settings.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *p.Settings.InitOptions.Lock))
	}

	if p.Settings.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", p.Settings.InitOptions.LockTimeout))
	}

	if !p.Settings.Refresh {
		args = append(args, "-refresh=false")
	}

	args = append(args, p.Settings.OutFile)

	cmd := execabs.Command(
		tofuBin,
		args...,
	)

	if !p.Settings.NoLog {
		trace.Cmd(cmd)
	}

	return cmd
}

func (p *Plugin) destroyCommand() *execabs.Cmd {
	args := []string{
		"destroy",
	}

	for _, v := range p.Settings.Targets.Value() {
		args = append(args, fmt.Sprintf("-target=%s", v))
	}

	if p.Settings.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", p.Settings.Parallelism))
	}

	if p.Settings.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *p.Settings.InitOptions.Lock))
	}

	if p.Settings.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", p.Settings.InitOptions.LockTimeout))
	}

	args = append(args, "-auto-approve")

	cmd := execabs.Command(
		tofuBin,
		args...,
	)

	if !p.Settings.NoLog {
		trace.Cmd(cmd)
	}

	return cmd
}
