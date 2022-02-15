package nebula

import (
	"github.com/glory-go/glory/log"
	nebula "github.com/vesoft-inc/nebula-go/v2"
)

var (
	logger nebula.Logger = &NebulaLogger{}
)

type NebulaLogger struct {
}

func (l *NebulaLogger) Info(msg string) {
	log.Info(msg)
}

func (l *NebulaLogger) Warn(msg string) {
	log.Warn(msg)
}

func (l *NebulaLogger) Error(msg string) {
	log.Error(msg)
}

func (l *NebulaLogger) Fatal(msg string) {
	log.Error(msg)
}
