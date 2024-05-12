package tofu

import (
	"fmt"

	"github.com/thegeeklab/wp-plugin-go/v2/types"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"
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

func (t *Tofu) Version() *types.Cmd {
	return &types.Cmd{
		Cmd:     execabs.Command(TofuBin, "version"),
		Private: t.NoLog,
	}
}

func (t *Tofu) Init() *types.Cmd {
	args := []string{
		"init",
	}

	for _, v := range t.InitOptions.BackendConfig {
		args = append(args, fmt.Sprintf("-backend-config=%s", v))
	}

	// Fail tofu execution on prompt
	args = append(args, "-input=false")

	return &types.Cmd{
		Cmd: execabs.Command(TofuBin, args...),
	}
}

func (t *Tofu) GetModules() *types.Cmd {
	return &types.Cmd{
		Cmd: execabs.Command(TofuBin, "get"),
	}
}

func (t *Tofu) Validate() *types.Cmd {
	return &types.Cmd{
		Cmd: execabs.Command(TofuBin, "validate"),
	}
}

func (t *Tofu) Fmt() *types.Cmd {
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

	return &types.Cmd{
		Cmd: execabs.Command(TofuBin, args...),
	}
}

func (t *Tofu) Plan(destroy bool) *types.Cmd {
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

	return &types.Cmd{
		Cmd:     execabs.Command(TofuBin, args...),
		Private: t.NoLog,
	}
}

func (t *Tofu) Apply() *types.Cmd {
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

	return &types.Cmd{
		Cmd:     execabs.Command(TofuBin, args...),
		Private: t.NoLog,
	}
}

func (t *Tofu) Destroy() *types.Cmd {
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

	return &types.Cmd{
		Cmd:     execabs.Command(TofuBin, args...),
		Private: t.NoLog,
	}
}
