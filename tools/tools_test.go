package tools

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	type args struct {
		envStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.envStr); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFromEnvIfNeed(t *testing.T) {
	type MetricsConfig struct {
		ConfigSource string `yaml:"config_source"`
		MetricsType  string `yaml:"metrics_type"`
		ActionType   string `yaml:"action_type"`
		ClientPort   string `yaml:"client_port"`
		ClientPath   string `yaml:"client_path"`
		GateWayHost  string `yaml:"gateway_host"`
		GateWayPort  string `yaml:"gateway_port"`
		JobName      string `yaml:"job_name"` // 数据上报job_name，默认为配置文件中的server_name
	}

	type args struct {
		rawConfig MetricsConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{rawConfig: MetricsConfig{
			ConfigSource: "env",
			MetricsType:  "1",
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadFromEnvIfNeed(&tt.args.rawConfig); (err != nil) != tt.wantErr {
				t.Errorf("ReadAllConfFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
