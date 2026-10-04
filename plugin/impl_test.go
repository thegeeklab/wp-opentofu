package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thegeeklab/wp-opentofu/tofu"
)

func TestFlagsFromContext(t *testing.T) {
	check := true

	tests := []struct {
		name     string
		initOpt  string
		fmtOpt   string
		wantInit tofu.InitOptions
		wantFmt  tofu.FmtOptions
		wantErr  bool
	}{
		{
			name:    "valid init and fmt options",
			initOpt: `{"backend-config":["foo=bar"],"lockfile":".terraform.lock.hcl"}`,
			fmtOpt:  `{"check":true}`,
			wantInit: tofu.InitOptions{
				BackendConfig: []string{"foo=bar"},
				Lockfile:      ".terraform.lock.hcl",
			},
			wantFmt: tofu.FmtOptions{Check: &check},
		},
		{
			name: "empty options are a no-op",
		},
		{
			name:    "invalid init option json",
			initOpt: `{"backend-config": invalid}`,
			wantErr: true,
		},
		{
			name:    "invalid fmt option json",
			fmtOpt:  `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PLUGIN_INIT_OPTION", tt.initOpt)
			t.Setenv("PLUGIN_FMT_OPTION", tt.fmtOpt)

			p := New(func(_ context.Context) error { return nil })
			_ = p.App.Run(t.Context(), []string{"wp-opentofu"})

			err := p.FlagsFromContext()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantInit, p.Settings.Tofu.InitOptions)
			assert.Equal(t, tt.wantFmt, p.Settings.Tofu.FmtOptions)
		})
	}
}
