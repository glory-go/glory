package configmanager

import "testing"

func TestGetConfigCenterPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigCenterPath(); got != tt.want {
				t.Errorf("GetConfigCenterPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
