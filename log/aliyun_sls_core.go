package log

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"

	"github.com/gogo/protobuf/proto"

	"github.com/google/martian/log"

	"go.uber.org/zap/zapcore"
)

import (
	"github.com/glory-go/glory/config"
)

type aliyunSLSCore struct {
	zapcore.LevelEnabler
	enc zapcore.Encoder
	out zapcore.WriteSyncer

	client       sls.ClientInterface
	projectName  string
	logstoreName string

	logGroup   *sls.LogGroup
	orgName    string
	serverName string

	lock               *sync.RWMutex
	tickerTimeInterval time.Duration
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func (c *aliyunSLSCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields)
	return clone
}

func (c *aliyunSLSCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *aliyunSLSCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	content := make([]*sls.LogContent, 0)
	content = append(content, &sls.LogContent{
		Key:   proto.String("level"),
		Value: proto.String(ent.Level.String()),
	})
	content = append(content, &sls.LogContent{
		Key:   proto.String("message"),
		Value: proto.String(ent.Message),
	})

	content = append(content, &sls.LogContent{
		Key:   proto.String("org"),
		Value: proto.String(c.orgName),
	})
	content = append(content, &sls.LogContent{
		Key:   proto.String("app"),
		Value: proto.String(c.serverName),
	})
	gloryE, ok := c.enc.(*GloryEncoder)
	if ok {
		content = append(content, &sls.LogContent{
			Key:   proto.String("logid"),
			Value: proto.String(gloryE.traceID),
		})
	}
	if ent.Level >= zapcore.ErrorLevel {
		content = append(content, &sls.LogContent{
			Key:   proto.String("stack"),
			Value: proto.String(ent.Stack),
		})
	}
	// add fields
	for _, field := range fields {
		content = append(content, &sls.LogContent{
			Key:   proto.String(field.Key),
			Value: proto.String(field.String),
		})
	}

	log := &sls.Log{
		Time:     proto.Uint32(uint32(ent.Time.Unix())),
		Contents: content,
	}

	c.lock.Lock()
	c.logGroup.Logs = append(c.logGroup.Logs, log)
	c.lock.Unlock()

	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		if err := c.Sync(); err != nil {
			return err
		}
	}
	return nil
}

func (c *aliyunSLSCore) Sync() error {
	var err error
	c.lock.Lock()
	defer c.lock.Unlock()
	for retryTimes := 0; retryTimes < 10; retryTimes++ {
		if len(c.logGroup.Logs) == 0 {
			break
		}
		err = c.client.PutLogs(c.projectName, c.logstoreName, c.logGroup)
		if err == nil {
			c.logGroup.Logs = make([]*sls.Log, 0)
			break
		} else {
			//handle exception here, you can add retryable erorrCode, set appropriate put_retry
			if strings.Contains(err.Error(), sls.WRITE_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.PROJECT_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.SHARD_WRITE_QUOTA_EXCEED) {
				//mayby you should split shard
				time.Sleep(1000 * time.Millisecond)
			} else if strings.Contains(err.Error(), sls.INTERNAL_SERVER_ERROR) || strings.Contains(err.Error(), sls.SERVER_BUSY) {
				time.Sleep(200 * time.Millisecond)
			} else {
				fmt.Printf("error: aliyun sls log sync failed with error = %v\n", err)
				break
			}
		}
	}
	return err
}

func (c *aliyunSLSCore) clone() *aliyunSLSCore {
	clone := *c
	clone.enc = c.enc.Clone()
	clone.out = c.out
	clone.lock = &sync.RWMutex{}
	return &clone
}

func (c *aliyunSLSCore) runUpload() {
	ticker := time.NewTicker(c.tickerTimeInterval)
	for range ticker.C {
		if err := c.Sync(); err != nil {
			log.Errorf("aliyun sls core updload sync failed with error = %s", err)
		}
	}
}

func newAliyunSLSLoggerCore(encoder zapcore.Encoder, enab zapcore.LevelEnabler, config *config.LogConfig, serverName, orgName string) *aliyunSLSCore {
	client := sls.CreateNormalInterface(config.EndPoint, config.AccessKeyID, config.AccessSecret, "")
	core := &aliyunSLSCore{
		enc:                encoder,
		client:             client,
		LevelEnabler:       enab,
		orgName:            orgName,
		serverName:         serverName,
		logstoreName:       config.LogStoreName,
		projectName:        config.ProjectName,
		logGroup:           &sls.LogGroup{},
		lock:               &sync.RWMutex{},
		tickerTimeInterval: time.Second * time.Duration(config.UploadInterval),
	}
	go core.runUpload()
	return core
}
