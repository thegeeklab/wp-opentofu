package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thegeeklab/wp-opentofu/tofu"
)

func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name  string
		input bool
		want  *bool
	}{
		{
			name:  "true value",
			input: true,
			want:  boolPtr(true),
		},
		{
			name:  "false value",
			input: false,
			want:  boolPtr(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boolPtr(tt.input)
			assert.Equal(t, tt.want, got)
			assert.NotNil(t, got)
			assert.Equal(t, tt.input, *got)
		})
	}
}

func TestParseStringMapToInitOptions(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  tofu.InitOptions
	}{
		{
			name:  "empty map",
			input: map[string]string{},
			want:  tofu.InitOptions{},
		},
		{
			name:  "backend true",
			input: map[string]string{"backend": "true"},
			want:  tofu.InitOptions{Backend: boolPtr(true)},
		},
		{
			name:  "backend false",
			input: map[string]string{"backend": "false"},
			want:  tofu.InitOptions{Backend: boolPtr(false)},
		},
		{
			name:  "backend-config",
			input: map[string]string{"backend-config": "config-value"},
			want:  tofu.InitOptions{BackendConfig: []string{"config-value"}},
		},
		{
			name:  "lock true",
			input: map[string]string{"lock": "true"},
			want:  tofu.InitOptions{Lock: boolPtr(true)},
		},
		{
			name:  "lockfile",
			input: map[string]string{"lockfile": "my.lock"},
			want:  tofu.InitOptions{Lockfile: "my.lock"},
		},
		{
			name:  "lock-timeout",
			input: map[string]string{"lock-timeout": "30s"},
			want:  tofu.InitOptions{LockTimeout: "30s"},
		},
		{
			name: "multiple options",
			input: map[string]string{
				"backend":  "true",
				"lock":     "false",
				"lockfile": "test.lock",
			},
			want: tofu.InitOptions{
				Backend:  boolPtr(true),
				Lock:     boolPtr(false),
				Lockfile: "test.lock",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseStringMapToInitOptions(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseStringMapToFmtOptions(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  tofu.FmtOptions
	}{
		{
			name:  "empty map",
			input: map[string]string{},
			want:  tofu.FmtOptions{},
		},
		{
			name:  "list true",
			input: map[string]string{"list": "true"},
			want:  tofu.FmtOptions{List: boolPtr(true)},
		},
		{
			name:  "list false",
			input: map[string]string{"list": "false"},
			want:  tofu.FmtOptions{List: boolPtr(false)},
		},
		{
			name:  "write true",
			input: map[string]string{"write": "true"},
			want:  tofu.FmtOptions{Write: boolPtr(true)},
		},
		{
			name:  "diff true",
			input: map[string]string{"diff": "true"},
			want:  tofu.FmtOptions{Diff: boolPtr(true)},
		},
		{
			name:  "check true",
			input: map[string]string{"check": "true"},
			want:  tofu.FmtOptions{Check: boolPtr(true)},
		},
		{
			name: "multiple options",
			input: map[string]string{
				"list":  "true",
				"write": "false",
				"diff":  "true",
				"check": "false",
			},
			want: tofu.FmtOptions{
				List:  boolPtr(true),
				Write: boolPtr(false),
				Diff:  boolPtr(true),
				Check: boolPtr(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseStringMapToFmtOptions(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
