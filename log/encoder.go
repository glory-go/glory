package log

import "go.uber.org/zap/zapcore"

type GloryEncoder struct {
	zapcore.Encoder

	traceID string
}

func NewGloryEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &GloryEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}

func (e *GloryEncoder) AddString(key, value string) {
	if key == GetTraceIDKey() {
		e.traceID = value
	} else {
		e.traceID = ""
		e.AddString(key, value)
	}
}

func (e *GloryEncoder) Clone() zapcore.Encoder {
	return &GloryEncoder{
		Encoder: e.Encoder.Clone(),
		traceID: e.traceID,
	}
}
