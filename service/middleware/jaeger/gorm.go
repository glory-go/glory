package jaeger

import (
	"context"
	"strings"
)

import (
	"github.com/opentracing/opentracing-go"

	"gorm.io/gorm"

	gormopentracing "gorm.io/plugin/opentracing"
)

type GormTracer struct {
	tracer opentracing.Tracer
}

func GormUseTrace(db *gorm.DB) error {
	return db.Use(gormopentracing.New(
		gormopentracing.WithTracer(tracer),
		gormopentracing.WithLogResult(false),
		gormopentracing.WithSqlParameters(true),
	))
}

// operationStage indicates the timing when the operation happens.
type operationStage string

// Name returns the actual string of operationStage.
func (op operationStage) Name() string {
	return string(op)
}

// operationName defines a type to wrap the name of each operation name.
type operationName string

// String returns the actual string of operationName.
func (op operationName) String() string {
	return string(op)
}

const (
	_stageBeforeCreate operationStage = "opentracing:before_create"
	_stageAfterCreate  operationStage = "opentracing:after_create"
	_stageBeforeUpdate operationStage = "opentracing:before_update"
	_stageAfterUpdate  operationStage = "opentracing:after_update"
	_stageBeforeQuery  operationStage = "opentracing:before_query"
	_stageAfterQuery   operationStage = "opentracing:after_query"
	_stageBeforeDelete operationStage = "opentracing:before_delete"
	_stageAfterDelete  operationStage = "opentracing:after_delete"
	_stageBeforeRow    operationStage = "opentracing:before_row"
	_stageAfterRow     operationStage = "opentracing:after_row"
	_stageBeforeRaw    operationStage = "opentracing:before_raw"
	_stageAfterRaw     operationStage = "opentracing:after_raw"
)

const (
	_createOp operationName = "create"
	_updateOp operationName = "update"
	_queryOp  operationName = "query"
	_deleteOp operationName = "delete"
	_rowOp    operationName = "row"
	_rawOp    operationName = "raw"
)

const (
	_prefix = "gorm.opentracing"
)

var (
	opentracingSpanKey = "opentracing:span"
)

var (
	// span.Tag keys
	_tableTagKey = keyWithPrefix("table")
)

func keyWithPrefix(key string) string {
	return _prefix + "." + key
}

func (t *GormTracer) Name() string {
	return "glory-gorm-tracer"
}

func (t *GormTracer) Initialize(db *gorm.DB) (err error) {
	e := myError{
		errs: make([]string, 0, 12),
	}

	// create
	err = db.Callback().Create().Before("gorm:create").Register(_stageBeforeCreate.Name(), t.beforeCreate)
	e.add(_stageBeforeCreate, err)
	err = db.Callback().Create().After("gorm:create").Register(_stageAfterCreate.Name(), t.after)
	e.add(_stageAfterCreate, err)

	// update
	err = db.Callback().Update().Before("gorm:update").Register(_stageBeforeUpdate.Name(), t.beforeUpdate)
	e.add(_stageBeforeUpdate, err)
	err = db.Callback().Update().After("gorm:update").Register(_stageAfterUpdate.Name(), t.after)
	e.add(_stageAfterUpdate, err)

	// query
	err = db.Callback().Query().Before("gorm:query").Register(_stageBeforeQuery.Name(), t.beforeQuery)
	e.add(_stageBeforeQuery, err)
	err = db.Callback().Query().After("gorm:query").Register(_stageAfterQuery.Name(), t.after)
	e.add(_stageAfterQuery, err)

	// delete
	err = db.Callback().Delete().Before("gorm:delete").Register(_stageBeforeDelete.Name(), t.beforeDelete)
	e.add(_stageBeforeDelete, err)
	err = db.Callback().Delete().After("gorm:delete").Register(_stageAfterDelete.Name(), t.after)
	e.add(_stageAfterDelete, err)

	// row
	err = db.Callback().Row().Before("gorm:row").Register(_stageBeforeRow.Name(), t.beforeRow)
	e.add(_stageBeforeRow, err)
	err = db.Callback().Row().After("gorm:row").Register(_stageAfterRow.Name(), t.after)
	e.add(_stageAfterRow, err)

	// raw
	err = db.Callback().Raw().Before("gorm:raw").Register(_stageBeforeRaw.Name(), t.beforeRaw)
	e.add(_stageBeforeRaw, err)
	err = db.Callback().Raw().After("gorm:raw").Register(_stageAfterRaw.Name(), t.after)
	e.add(_stageAfterRaw, err)

	return e.toError()
}

func (t *GormTracer) beforeCreate(db *gorm.DB) {
	t.injectBefore(db, _createOp)
}

func (t *GormTracer) after(db *gorm.DB) {
	t.extractAfter(db)
}

func (t *GormTracer) beforeUpdate(db *gorm.DB) {
	t.injectBefore(db, _updateOp)
}

func (t *GormTracer) beforeQuery(db *gorm.DB) {
	t.injectBefore(db, _queryOp)
}

func (t *GormTracer) beforeDelete(db *gorm.DB) {
	t.injectBefore(db, _deleteOp)
}

func (t *GormTracer) beforeRow(db *gorm.DB) {
	t.injectBefore(db, _rowOp)
}

func (t *GormTracer) beforeRaw(db *gorm.DB) {
	t.injectBefore(db, _rawOp)
}

func (t *GormTracer) injectBefore(db *gorm.DB, op operationName) {
	// make sure context could be used
	if db == nil {
		return
	}

	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "could not inject sp from nil Statement.Context or nil Statement")
		return
	}

	sp, _ := opentracing.StartSpanFromContextWithTracer(db.Statement.Context, t.tracer, op.String())
	db.InstanceSet(opentracingSpanKey, sp)
}

func (t *GormTracer) extractAfter(db *gorm.DB) {
	// make sure context could be used
	if db == nil {
		return
	}
	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "could not extract sp from nil Statement.Context or nil Statement")
		return
	}

	// extract sp from db context
	//sp := opentracing.SpanFromContext(db.Statement.Context)
	v, ok := db.InstanceGet(opentracingSpanKey)
	if !ok || v == nil {
		return
	}

	sp, ok := v.(opentracing.Span)
	if !ok || sp == nil {
		return
	}
	defer sp.Finish()

	// tag and log fields we want.
	tag(sp, db)
}

// tag called after operation
func tag(sp opentracing.Span, db *gorm.DB) {
	sp.SetTag(_tableTagKey, db.Statement.Table)
}

type myError struct {
	errs []string
}

func (e *myError) add(stage operationStage, err error) {
	if err == nil {
		return
	}

	e.errs = append(e.errs, "stage="+stage.Name()+":"+err.Error())
}

func (e myError) toError() error {
	if len(e.errs) == 0 {
		return nil
	}

	return e
}

func (e myError) Error() string {
	return strings.Join(e.errs, ";")
}
