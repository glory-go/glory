package file

import (
	logruscomponent "github.com/glory-go/glory/v2/components/log/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	FileHookType = "file"
)

func init() {
	logruscomponent.RegisterHookBuilder(FileHookType, hookBuilder)
}

func hookBuilder(raw map[string]any) (logrus.Hook, error) {
	conf := &config{}
	if err := mapstructure.Decode(raw, conf); err != nil {
		return nil, err
	}

	levelPath := lfshook.PathMap{}
	for level, path := range conf.LevelPath {
		parsedLevel, err := logrus.ParseLevel(level)
		if err != nil {
			return nil, err
		}
		levelPath[parsedLevel] = path
	}
	hook := lfshook.NewHook(conf.LevelPath, &logrus.JSONFormatter{})
	return hook, nil
}
