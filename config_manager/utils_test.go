package configmanager

import "testing"

func TestReadFromEnvIfNeed(t *testing.T) {
	type args struct {
		conf map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReadFromEnvIfNeed(tt.args.conf)
		})
	}
}
