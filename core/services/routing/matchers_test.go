package routing_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/services/routing"
)

const (
	maxContentLength = 2 << 16
)

func TestJsonPath(t *testing.T) {
	type args struct {
		path string
		want any
	}
	tests := []struct {
		name      string
		args      args
		req       *http.Request
		wantMatch bool
		matchErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Simple JSON patch",
			args: args{
				path: "$.name",
				want: "Ted",
			},
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"name":"Ted"}`)),
			},
			wantMatch: true,
			matchErr:  assert.NoError,
		},
		{
			name: "Multiple JSON matches",
			args: args{
				path: "[*].name",
				want: "Ted",
			},
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`[{"name":"Will"},{"name":"Ted"}]`)),
			},
			wantMatch: true,
			matchErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routing.JsonPath(tt.args.path, tt.args.want)
			if !tt.matchErr(t, err) {
				return
			}

			assert.Equal(
				t,
				tt.wantMatch,
				got.Matches(domain.NewRequest(tt.req)),
				"Matcher returned unexpected response",
			)
		})
	}
}
