package grammar_test

import (
	"testing"

	grammar2 "github.com/prskr/go-dito/core/services/grammar"
)

func TestParam_AsString(t *testing.T) {
	t.Parallel()
	type fields struct {
		String *string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Empty string",
			fields: fields{
				String: grammar2.StringP(""),
			},
			want: "",
		},
		{
			name: "Any string",
			fields: fields{
				String: grammar2.StringP("Hello, world!"),
			},
			want: "Hello, world!",
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := grammar2.Param{
				String: tt.fields.String,
			}
			got, err := p.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParam_AsInt(t *testing.T) {
	t.Parallel()
	type fields struct {
		Int *int
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "zero value",
			fields: fields{
				Int: grammar2.IntP(0),
			},
			want: 0,
		},
		{
			name: "Any int",
			fields: fields{
				Int: grammar2.IntP(42),
			},
			want: 42,
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := grammar2.Param{
				Int: tt.fields.Int,
			}
			got, err := p.AsInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
