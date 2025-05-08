package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestYamlToJson(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:  "empty",
			input: "---",
			want:  []byte("null"),
		},
		{
			name:  "simple string",
			input: `hello`,
			want:  []byte(`"hello"`),
		},
		{
			name:  "simple number",
			input: `42`,
			want:  []byte(`42`),
		},
		{
			name:  "simple float",
			input: `42.13`,
			want:  []byte(`42.13`),
		},
		{
			name: "Array of strings",
			input: `
- hello
- world
`,
			want: []byte(`["hello","world"]`),
		},
		{
			name: "map of strings",
			input: `
someKey:
  nestedKey: hello
  world: "ted!"
`,
			want: []byte(`{"someKey":{"nestedKey":"hello","world":"ted!"}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var node yaml.Node
			assert.NoError(t, yaml.Unmarshal([]byte(tt.input), &node))
			got, err := YamlToJson(&node)
			if !assert.NoError(t, err, "YamlToJson() error = %v, wantErr %v", err, tt.wantErr) {
				return
			}

			assert.Equal(t, string(tt.want), string(got))
		})
	}
}
