package plugin

import (
	"fmt"

	"golang.org/x/sys/execabs"
)

const (
	tofuBin = "/usr/local/bin/tofu"
)

func (p *Plugin) versionCommand() *Cmd {
	return &Cmd{
		execabs.Command(tofuBin, "version"),
		p.Settings.NoLog,
	}
}

func (p *Plugin) initCommand() *Cmd {
	args := []string{
		"init",
	}

	for _, v := range p.Settings.InitOptions.BackendConfig {
		args = append(args, fmt.Sprintf("-backend-config=%s", v))
	}

	// Fail tofu execution on prompt
	args = append(args, "-input=false")

	return &Cmd{
		execabs.Command(tofuBin, args...),
		false,
	}
}

func (p *Plugin) getModulesCommand() *Cmd {
	return &Cmd{
		execabs.Command(tofuBin, "get"),
		false,
	}
}

func (p *Plugin) validateCommand() *Cmd {
	return &Cmd{
		execabs.Command(tofuBin, "validate"),
		false,
	}
}

func (p *Plugin) fmtCommand() *Cmd {
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

	return &Cmd{
		execabs.Command(tofuBin, args...),
		false,
	}
}

func (p *Plugin) planCommand(destroy bool) *Cmd {
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

	return &Cmd{
		execabs.Command(tofuBin, args...),
		p.Settings.NoLog,
	}
}

func (p *Plugin) applyCommand() *Cmd {
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

	return &Cmd{
		execabs.Command(tofuBin, args...),
		p.Settings.NoLog,
	}
}

func (p *Plugin) destroyCommand() *Cmd {
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

	return &Cmd{
		execabs.Command(tofuBin, args...),
		p.Settings.NoLog,
	}
}
