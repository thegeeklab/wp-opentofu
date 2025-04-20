package plugin

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
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
			_ = got.FlagsFromContext()

			assert.ElementsMatch(t, tt.want, got.Environment.Value())
		})
	}
}
