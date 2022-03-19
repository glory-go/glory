package config

import (
	"reflect"
	"testing"
)

func init() {

}

func Test_replaceMapValueFromConfigCenter(t *testing.T) {
	type args struct {
		configSource string
		conf         map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replaceMapValueFromConfigCenter(tt.args.configSource, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceMapValueFromConfigCenter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replaceMapValueFromConfigCenter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceStringValueFromConfigCenter(t *testing.T) {
	type args struct {
		configSource string
		rawVal       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replaceStringValueFromConfigCenter(tt.args.configSource, tt.args.rawVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceStringValueFromConfigCenter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("replaceStringValueFromConfigCenter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGroupAndKey(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name      string
		args      args
		wantGroup string
		wantKey   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroup, gotKey := getGroupAndKey(tt.args.val)
			if gotGroup != tt.wantGroup {
				t.Errorf("getGroupAndKey() gotGroup = %v, want %v", gotGroup, tt.wantGroup)
			}
			if gotKey != tt.wantKey {
				t.Errorf("getGroupAndKey() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
