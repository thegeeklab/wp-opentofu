package tofu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestTofu_Version(t *testing.T) {
	tests := []struct {
		name string
		tofu *Tofu
		want []string
	}{
		{
			name: "test version command",
			tofu: &Tofu{},
			want: []string{TofuBin, "version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.Version()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestTofu_Init(t *testing.T) {
	tests := []struct {
		name string
		tofu *Tofu
		want []string
	}{
		{
			name: "init with no backend config",
			tofu: &Tofu{},
			want: []string{
				TofuBin,
				"init",
				"-input=false",
			},
		},
		{
			name: "init with single backend config",
			tofu: &Tofu{
				InitOptions: InitOptions{
					BackendConfig: []string{"key=value"},
				},
			},
			want: []string{
				TofuBin,
				"init",
				"-backend-config=key=value",
				"-input=false",
			},
		},
		{
			name: "init with multiple backend configs",
			tofu: &Tofu{
				InitOptions: InitOptions{
					BackendConfig: []string{"key1=value1", "key2=value2"},
				},
			},
			want: []string{
				TofuBin,
				"init",
				"-backend-config=key1=value1",
				"-backend-config=key2=value2",
				"-input=false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.Init()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestTofu_GetModules(t *testing.T) {
	tests := []struct {
		name string
		tofu *Tofu
		want []string
	}{
		{
			name: "get modules command",
			tofu: &Tofu{},
			want: []string{TofuBin, "get"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.GetModules()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestTofu_Validate(t *testing.T) {
	tests := []struct {
		name string
		tofu *Tofu
		want []string
	}{
		{
			name: "validate command",
			tofu: &Tofu{},
			want: []string{TofuBin, "validate"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.Validate()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestTofu_Fmt(t *testing.T) {
	tests := []struct {
		name string
		tofu *Tofu
		want []string
	}{
		{
			name: "fmt with no options",
			tofu: &Tofu{},
			want: []string{
				TofuBin,
				"fmt",
			},
		},
		{
			name: "fmt with list option",
			tofu: &Tofu{
				FmtOptions: FmtOptions{
					List: boolPtr(true),
				},
			},
			want: []string{
				TofuBin,
				"fmt",
				"-list=true",
			},
		},
		{
			name: "fmt with write option",
			tofu: &Tofu{
				FmtOptions: FmtOptions{
					Write: boolPtr(true),
				},
			},
			want: []string{
				TofuBin,
				"fmt",
				"-write=true",
			},
		},
		{
			name: "fmt with diff option",
			tofu: &Tofu{
				FmtOptions: FmtOptions{
					Diff: boolPtr(true),
				},
			},
			want: []string{
				TofuBin,
				"fmt",
				"-diff=true",
			},
		},
		{
			name: "fmt with check option",
			tofu: &Tofu{
				FmtOptions: FmtOptions{
					Check: boolPtr(true),
				},
			},
			want: []string{
				TofuBin,
				"fmt",
				"-check=true",
			},
		},
		{
			name: "fmt with multiple options",
			tofu: &Tofu{
				FmtOptions: FmtOptions{
					List:  boolPtr(true),
					Write: boolPtr(true),
					Diff:  boolPtr(true),
					Check: boolPtr(true),
				},
			},
			want: []string{
				TofuBin,
				"fmt",
				"-list=true",
				"-write=true",
				"-diff=true",
				"-check=true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.Fmt()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestTofu_Plan(t *testing.T) {
	tests := []struct {
		name    string
		tofu    *Tofu
		destroy bool
		want    []string
	}{
		{
			name:    "plan with no options",
			tofu:    &Tofu{},
			destroy: false,
			want: []string{
				TofuBin,
				"plan",
				"-refresh=false",
			},
		},
		{
			name: "plan with output options",
			tofu: &Tofu{
				OutFile: "plan.tfout",
			},
			destroy: false,
			want: []string{
				TofuBin,
				"plan",
				"-out=plan.tfout",
				"-refresh=false",
			},
		},
		{
			name:    "plan with destroy option",
			tofu:    &Tofu{},
			destroy: true,
			want: []string{
				TofuBin,
				"plan",
				"-destroy",
				"-refresh=false",
			},
		},
		{
			name: "plan with targets",
			tofu: &Tofu{
				Targets: *cli.NewStringSlice("target1", "target2"),
			},
			destroy: false,
			want: []string{
				TofuBin,
				"plan",
				"--target", "target1",
				"--target", "target2",
				"-refresh=false",
			},
		},
		{
			name: "plan with parallelism",
			tofu: &Tofu{
				Parallelism: 10,
			},
			destroy: false,
			want: []string{
				TofuBin,
				"plan",
				"-parallelism=10",
				"-refresh=false",
			},
		},
		{
			name: "plan with lock option",
			tofu: &Tofu{
				InitOptions: InitOptions{
					Lock: boolPtr(true),
				},
			},
			destroy: false,
			want: []string{
				TofuBin,
				"plan",
				"-lock=true",
				"-refresh=false",
			},
		},
		{
			name: "plan with lock timeout",
			tofu: &Tofu{
				InitOptions: InitOptions{
					LockTimeout: "10s",
				},
			},
			destroy: false,
			want: []string{
				TofuBin,
				"plan",
				"-lock-timeout=10s",
				"-refresh=false",
			},
		},
		{
			name: "plan with refresh option",
			tofu: &Tofu{
				Refresh: true,
			},
			destroy: false,
			want: []string{
				TofuBin,
				"plan",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.Plan(tt.destroy)
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestTofu_Apply(t *testing.T) {
	tests := []struct {
		name string
		tofu *Tofu
		want []string
	}{
		{
			name: "apply with no options",
			tofu: &Tofu{},
			want: []string{
				TofuBin,
				"apply",
				"-refresh=false",
			},
		},
		{
			name: "apply with targets",
			tofu: &Tofu{
				Targets: *cli.NewStringSlice("target1", "target2"),
			},
			want: []string{
				TofuBin,
				"apply",
				"--target", "target1",
				"--target", "target2",
				"-refresh=false",
			},
		},
		{
			name: "apply with parallelism",
			tofu: &Tofu{
				Parallelism: 10,
			},
			want: []string{
				TofuBin,
				"apply",
				"-parallelism=10",
				"-refresh=false",
			},
		},
		{
			name: "apply with lock option",
			tofu: &Tofu{
				InitOptions: InitOptions{
					Lock: boolPtr(true),
				},
			},
			want: []string{
				TofuBin,
				"apply",
				"-lock=true",
				"-refresh=false",
			},
		},
		{
			name: "apply with lock timeout",
			tofu: &Tofu{
				InitOptions: InitOptions{
					LockTimeout: "10s",
				},
			},
			want: []string{
				TofuBin,
				"apply",
				"-lock-timeout=10s",
				"-refresh=false",
			},
		},
		{
			name: "apply with refresh option",
			tofu: &Tofu{
				Refresh: true,
			},
			want: []string{
				TofuBin,
				"apply",
			},
		},
		{
			name: "apply with output file",
			tofu: &Tofu{
				OutFile: "out.tfout",
			},
			want: []string{
				TofuBin,
				"apply",
				"-refresh=false",
				"out.tfout",
			},
		},
		{
			name: "apply with no log",
			tofu: &Tofu{
				NoLog: true,
			},
			want: []string{
				TofuBin,
				"apply",
				"-refresh=false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.Apply()
			assert.Equal(t, tt.want, cmd.Args)

			if tt.tofu.NoLog {
				assert.Equal(t, cmd.Stdout, nil)
			}
		})
	}
}

func TestTofu_Destroy(t *testing.T) {
	tests := []struct {
		name string
		tofu *Tofu
		want []string
	}{
		{
			name: "destroy with no options",
			tofu: &Tofu{},
			want: []string{
				TofuBin,
				"destroy",
				"-auto-approve",
			},
		},
		{
			name: "destroy with targets",
			tofu: &Tofu{
				Targets: *cli.NewStringSlice("target1", "target2"),
			},
			want: []string{
				TofuBin,
				"destroy",
				"-target=target1",
				"-target=target2",
				"-auto-approve",
			},
		},
		{
			name: "destroy with parallelism",
			tofu: &Tofu{
				Parallelism: 10,
			},
			want: []string{
				TofuBin,
				"destroy",
				"-parallelism=10",
				"-auto-approve",
			},
		},
		{
			name: "destroy with lock option",
			tofu: &Tofu{
				InitOptions: InitOptions{
					Lock: boolPtr(true),
				},
			},
			want: []string{
				TofuBin,
				"destroy",
				"-lock=true",
				"-auto-approve",
			},
		},
		{
			name: "destroy with lock timeout",
			tofu: &Tofu{
				InitOptions: InitOptions{
					LockTimeout: "10s",
				},
			},
			want: []string{
				TofuBin,
				"destroy",
				"-lock-timeout=10s",
				"-auto-approve",
			},
		},
		{
			name: "destroy with no log",
			tofu: &Tofu{
				NoLog: true,
			},
			want: []string{
				TofuBin,
				"destroy",
				"-auto-approve",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.tofu.Destroy()

			assert.Equal(t, tt.want, cmd.Args)

			if tt.tofu.NoLog {
				assert.Equal(t, cmd.Stdout, nil)
			}
		})
	}
}
