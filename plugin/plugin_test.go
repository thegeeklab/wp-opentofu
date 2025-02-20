package plugin

import (
	"context"
	"reflect"
	"testing"
)

func TestFlags(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		want []string
	}{
		{
			name: "parse secrets list with escape",
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
		for key, value := range tt.envs {
			t.Setenv(key, value)
		}

		got := New(func(_ context.Context) error { return nil })

		_ = got.App.Run([]string{"wp-opentofu"})
		_ = got.FlagsFromContext()

		if !reflect.DeepEqual(got.Plugin.Environment.Value(), tt.want) {
			t.Errorf("%q. Plugin.Environment = %v, want %v", tt.name, got.Plugin.Environment.Value(), tt.want)
		}
	}
}
