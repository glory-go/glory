package log

import (
	"context"
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/olivere/elastic/v7"

	"go.uber.org/zap/zapcore"
)

import (
	"github.com/glory-go/glory/tools"
)

var (
	esLogIndexStruct = map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"service_name": map[string]interface{}{
					"type": "text",
				},
				"level": map[string]interface{}{
					"type": "integer",
				},
				"log": map[string]interface{}{
					"type": "text",
				},
				"@timestamp": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
	LogIndexName = "go-online"
)

func addData(client *elastic.Client, data interface{}) error {
	go func() {
		// TODO: 更改成业务的ctx
		ctx := context.Background()
		exists, err := client.IndexExists(LogIndexName).Do(ctx)
		if err != nil {
			fmt.Println("error push to elastic err:", err)
		}
		if !exists {
			createIndex, err := client.CreateIndex(LogIndexName).BodyJson(esLogIndexStruct).Do(ctx)
			if err != nil {
				fmt.Println("error push to elastic err:", err)
			}
			if !createIndex.Acknowledged {
				// Not acknowledged
				fmt.Println("error push to elastic, ack not received")
			}
		}
		if _, err = client.Index().
			Index(LogIndexName).
			BodyJson(data).
			Id(tools.GenerateXID()).
			Do(ctx); err != nil {
			fmt.Println("error push to elastic err:", err)
		}
	}()
	return nil
}

type ElasticDebugHook struct {
	es          *elastic.Client
	mux         sync.Mutex
	serviceName string
}

func (c *ElasticDebugHook) Write(p []byte) (n int, err error) {
	defer func() {
		c.mux.Unlock()
	}()
	item := map[string]interface{}{
		"service_name": c.serviceName,
		"@timestamp":   time.Now().Unix(),
		"level":        int(zapcore.DebugLevel),
		"log":          string(p),
	}
	c.mux.Lock()
	if err := addData(c.es, item); err != nil {
		return 0, err
	}
	return len(p), nil
}

type ElasticInfoHook struct {
	es          *elastic.Client
	mux         sync.Mutex
	serviceName string
}

func (c *ElasticInfoHook) Write(p []byte) (n int, err error) {
	defer func() {
		c.mux.Unlock()
	}()
	item := map[string]interface{}{
		"service_name": c.serviceName,
		"@timestamp":   time.Now().Unix(),
		"level":        int(zapcore.DebugLevel),
		"log":          string(p),
	}
	c.mux.Lock()
	if err := addData(c.es, item); err != nil {
		return 0, err
	}
	return len(p), nil
}

type ElasticWarnHook struct {
	es          *elastic.Client
	mux         sync.Mutex
	serviceName string
}

func (c *ElasticWarnHook) Write(p []byte) (n int, err error) {
	defer func() {
		c.mux.Unlock()
	}()
	item := map[string]interface{}{
		"service_name": c.serviceName,
		"@timestamp":   time.Now().Unix(),
		"level":        int(zapcore.DebugLevel),
		"log":          string(p),
	}
	c.mux.Lock()
	if err := addData(c.es, item); err != nil {
		return 0, err
	}
	return len(p), nil
}

type ElasticErrorHook struct {
	es          *elastic.Client
	mux         sync.Mutex
	serviceName string
}

func (c *ElasticErrorHook) Write(p []byte) (n int, err error) {
	defer func() {
		c.mux.Unlock()
	}()
	item := map[string]interface{}{
		"service_name": c.serviceName,
		"@timestamp":   time.Now().Unix(),
		"level":        int(zapcore.DebugLevel),
		"log":          string(p),
	}
	c.mux.Lock()
	if err := addData(c.es, item); err != nil {
		return 0, err
	}
	return len(p), nil
}

type ElasticPanicHook struct {
	es          *elastic.Client
	mux         sync.Mutex
	serviceName string
}

func (c *ElasticPanicHook) Write(p []byte) (n int, err error) {
	defer func() {
		c.mux.Unlock()
	}()
	item := map[string]interface{}{
		"service_name": c.serviceName,
		"@timestamp":   time.Now().Unix(),
		"level":        int(zapcore.DebugLevel),
		"log":          string(p),
	}
	c.mux.Lock()
	if err := addData(c.es, item); err != nil {
		return 0, err
	}
	return len(p), nil
}
