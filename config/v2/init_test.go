package config

import "testing"

func TestGetConfigPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigPath(); got != tt.want {
				t.Errorf("GetConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
