package logrus

import "github.com/sirupsen/logrus"

type logrusComponentConfig struct {
	Level logrus.Level `mapstructure:"level"`
}
