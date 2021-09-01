package jaeger

import (
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

func GormUseTrace(db *gorm.DB) error {
	return db.Use(gormopentracing.New(
		gormopentracing.WithTracer(tracer),
		gormopentracing.WithLogResult(false),
		gormopentracing.WithSqlParameters(true),
	))
}
