package plugin

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thegeeklab/wp-opentofu/tofu"
	"github.com/urfave/cli/v3"
)

func setupPluginTest(t *testing.T) *Plugin {
	t.Helper()

	cli.HelpPrinter = func(_ io.Writer, _ string, _ interface{}) {}
	got := New(func(_ context.Context) error { return nil })
	_ = got.App.Run(t.Context(), []string{"wp-docker-buildx"})

	return got
}

func TestEnvironmentFlag(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		want []string
	}{
		{
			name: "simple environment",
			envs: map[string]string{
				"PLUGIN_ENVIRONMENT": `{"env1": "value1", "env2": "value2"}`,
			},
			want: []string{
				"env1=value1",
				"env2=value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envs {
				t.Setenv(key, value)
			}

			got := setupPluginTest(t)
			_ = got.Validate()

			assert.ElementsMatch(t, tt.want, got.Environment.Value())
		})
	}
}

func TestOptionFlags(t *testing.T) {
	tests := []struct {
		name            string
		envs            map[string]string
		wantInitOptions tofu.InitOptions
		wantFmtOptions  tofu.FmtOptions
	}{
		{
			name: "init options parsing",
			envs: map[string]string{
				"PLUGIN_INIT_OPTION": `{"backend":"true","lock":"false","lockfile":"test.lock"}`,
			},
			wantInitOptions: tofu.InitOptions{
				Backend:  boolPtr(true),
				Lock:     boolPtr(false),
				Lockfile: "test.lock",
			},
			wantFmtOptions: tofu.FmtOptions{},
		},
		{
			name: "fmt options parsing",
			envs: map[string]string{
				"PLUGIN_FMT_OPTION": `{"list":"true","write":"false","diff":"true"}`,
			},
			wantInitOptions: tofu.InitOptions{},
			wantFmtOptions: tofu.FmtOptions{
				List:  boolPtr(true),
				Write: boolPtr(false),
				Diff:  boolPtr(true),
			},
		},
		{
			name: "both init and fmt options",
			envs: map[string]string{
				"PLUGIN_INIT_OPTION": `{"backend":"true","backend-config":"config-value"}`,
				"PLUGIN_FMT_OPTION":  `{"check":"true","write":"false"}`,
			},
			wantInitOptions: tofu.InitOptions{
				Backend:       boolPtr(true),
				BackendConfig: []string{"config-value"},
			},
			wantFmtOptions: tofu.FmtOptions{
				Check: boolPtr(true),
				Write: boolPtr(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envs {
				t.Setenv(key, value)
			}

			got := setupPluginTest(t)
			err := got.Validate()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantInitOptions, got.Settings.Tofu.InitOptions)
			assert.Equal(t, tt.wantFmtOptions, got.Settings.Tofu.FmtOptions)
		})
	}
}
