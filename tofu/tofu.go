package tofu

import (
	"fmt"
	"os"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v3/exec"
	"github.com/urfave/cli/v2"
)

const TofuBin = "/usr/local/bin/tofu"

type Tofu struct {
	InitOptions InitOptions
	FmtOptions  FmtOptions

	OutFile     string
	Parallelism int
	Targets     cli.StringSlice
	Refresh     bool
	NoLog       bool
}

// InitOptions include options for the OpenTofu init command.
type InitOptions struct {
	Backend       *bool    `json:"backend"`
	BackendConfig []string `json:"backend-config"`
	Lock          *bool    `json:"lock"`
	LockTimeout   string   `json:"lock-timeout"`
	Lockfile      string   `json:"lockfile"`
}

// FmtOptions fmt options for the OpenTofu fmt command.
type FmtOptions struct {
	List  *bool `json:"list"`
	Write *bool `json:"write"`
	Diff  *bool `json:"diff"`
	Check *bool `json:"check"`
}

func (t *Tofu) Version() *plugin_exec.Cmd {
	cmd := plugin_exec.Command(TofuBin, "version")

	if !t.NoLog {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func (t *Tofu) Init() *plugin_exec.Cmd {
	args := []string{
		"init",
	}

	if t.InitOptions.Backend != nil {
		args = append(args, fmt.Sprintf("-backend=%t", *t.InitOptions.Backend))
	}

	for _, v := range t.InitOptions.BackendConfig {
		args = append(args, fmt.Sprintf("-backend-config=%s", v))
	}

	if t.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *t.InitOptions.Lock))
	}

	if t.InitOptions.Lockfile != "" {
		args = append(args, fmt.Sprintf("-lockfile=%s", t.InitOptions.Lockfile))
	}

	if t.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", t.InitOptions.LockTimeout))
	}

	// Fail tofu execution on prompt
	args = append(args, "-input=false")

	cmd := plugin_exec.Command(TofuBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (t *Tofu) GetModules() *plugin_exec.Cmd {
	cmd := plugin_exec.Command(TofuBin, "get")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (t *Tofu) Validate() *plugin_exec.Cmd {
	cmd := plugin_exec.Command(TofuBin, "validate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (t *Tofu) Fmt() *plugin_exec.Cmd {
	args := []string{
		"fmt",
	}

	if t.FmtOptions.List != nil {
		args = append(args, fmt.Sprintf("-list=%t", *t.FmtOptions.List))
	}

	if t.FmtOptions.Write != nil {
		args = append(args, fmt.Sprintf("-write=%t", *t.FmtOptions.Write))
	}

	if t.FmtOptions.Diff != nil {
		args = append(args, fmt.Sprintf("-diff=%t", *t.FmtOptions.Diff))
	}

	if t.FmtOptions.Check != nil {
		args = append(args, fmt.Sprintf("-check=%t", *t.FmtOptions.Check))
	}

	cmd := plugin_exec.Command(TofuBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (t *Tofu) Plan(destroy bool) *plugin_exec.Cmd {
	args := []string{
		"plan",
	}

	if destroy {
		args = append(args, "-destroy")
	} else if t.OutFile != "" {
		args = append(args, fmt.Sprintf("-out=%s", t.OutFile))
	}

	for _, value := range t.Targets.Value() {
		args = append(args, "--target", value)
	}

	if t.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", t.Parallelism))
	}

	if t.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *t.InitOptions.Lock))
	}

	if t.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", t.InitOptions.LockTimeout))
	}

	if !t.Refresh {
		args = append(args, "-refresh=false")
	}

	cmd := plugin_exec.Command(TofuBin, args...)

	if !t.NoLog {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func (t *Tofu) Apply() *plugin_exec.Cmd {
	args := []string{
		"apply",
	}

	for _, v := range t.Targets.Value() {
		args = append(args, "--target", v)
	}

	if t.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", t.Parallelism))
	}

	if t.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *t.InitOptions.Lock))
	}

	if t.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", t.InitOptions.LockTimeout))
	}

	if !t.Refresh {
		args = append(args, "-refresh=false")
	}

	if t.OutFile != "" {
		args = append(args, t.OutFile)
	}

	cmd := plugin_exec.Command(TofuBin, args...)

	if !t.NoLog {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func (t *Tofu) Destroy() *plugin_exec.Cmd {
	args := []string{
		"destroy",
	}

	for _, v := range t.Targets.Value() {
		args = append(args, fmt.Sprintf("-target=%s", v))
	}

	if t.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", t.Parallelism))
	}

	if t.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *t.InitOptions.Lock))
	}

	if t.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", t.InitOptions.LockTimeout))
	}

	args = append(args, "-auto-approve")

	cmd := plugin_exec.Command(TofuBin, args...)

	if !t.NoLog {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
