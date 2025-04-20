package config_test

import (
	"testing"

	"github.com/prskr/go-dito/core/services/config"
)

func TestParseDataSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		rawDataSize string
		want        config.DataSize
		wantErr     bool
	}{
		{
			name:        "Plain number",
			rawDataSize: "100",
			want:        config.DataSize(100),
			wantErr:     false,
		},
		{
			name:        "With bytes suffix",
			rawDataSize: "100b",
			want:        config.DataSize(100),
			wantErr:     false,
		},
		{
			name:        "With kb suffix",
			rawDataSize: "100kb",
			want:        100 * config.KiloByte,
			wantErr:     false,
		},
		{
			name:        "With mb suffix",
			rawDataSize: "5mb",
			want:        5 * config.MegaByte,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := config.ParseDataSize(tt.rawDataSize)

			if (err != nil) != tt.wantErr {
				t.Errorf("Want error %t but got %v", tt.wantErr, err)
			}

			if err != nil {
				return
			}

			if got != tt.want {
				t.Errorf("Got %d but expected %d", got, tt.want)
			}
		})
	}
}

func TestDataSize_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rawJson []byte
		want    config.DataSize
		wantErr bool
	}{
		{
			name:    "Plain number",
			rawJson: []byte(`100`),
			want:    config.DataSize(100),
			wantErr: false,
		},
		{
			name:    "Number with unit",
			rawJson: []byte(`100kb`),
			want:    100 * config.KiloByte,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var value config.DataSize
			if err := value.UnmarshalJSON(tt.rawJson); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if value != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", value, tt.want)
			}
		})
	}
}
