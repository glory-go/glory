package file

import "github.com/rifflock/lfshook"

type config struct {
	LevelPath lfshook.PathMap `mapstructure:"level_path"`
}
