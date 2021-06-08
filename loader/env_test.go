package loader

import (
	"os"
	"testing"

	"github.com/glory-go/glory/config"
	"github.com/stretchr/testify/assert"
)

type testConf struct {
	A string
	B []string
	C map[string]string
	D struct {
		D1 string
		D2 []string
		D3 int
		D4 map[string]int
	}
	// E map[string]struct {
	// 	E1 string
	// 	E2 []string
	// }
}

func (*testConf) Name() string {
	return "test"
}
func (*testConf) ConfigSource() config.ConfigSource {
	return EnvLoaderName
}
func (*testConf) Init() error {
	return nil
}

func getConf(a, b string) config.ConfigInterface {
	return &testConf{
		A: a,
		B: []string{a, b},
		C: map[string]string{
			"a": a,
			"b": b,
		},
		D: struct {
			D1 string
			D2 []string
			D3 int
			D4 map[string]int
		}{
			D1: a,
			D2: []string{a, b},
			D3: 1,
			D4: map[string]int{
				"a": 1,
				"b": 1,
			},
		},
		// E: map[string]struct {
		// 	E1 string
		// 	E2 []string
		// }{
		// 	"a": {
		// 		E1: a,
		// 		E2: []string{a, b},
		// 	},
		// },
	}
}

func TestEnvLoader_Load(t *testing.T) {
	type args struct {
		config config.ConfigInterface
	}
	tt := struct {
		name    string
		e       EnvLoader
		args    args
		want    config.ConfigInterface
		wantErr bool
	}{
		name: "复合结构",
		e:    EnvLoader{},
		args: args{
			config: getConf("a", "b"),
		},
		want:    getConf("b", "b"),
		wantErr: false,
	}
	t.Run(tt.name, func(t *testing.T) {
		os.Setenv("a", "b")
		os.Setenv("b", "")
		if err := tt.e.Load(tt.args.config); (err != nil) != tt.wantErr {
			t.Errorf("EnvLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
		}
		t.Log(tt.args.config, tt.want)
		assert.True(t, assert.ObjectsAreEqual(tt.want, tt.args.config))
	})
}
