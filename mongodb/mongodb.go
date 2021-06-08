package mongodb

import (
	"errors"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBHandler struct {
	services map[string]*MongoDBService
}

var defaultMongoDBHandler *MongoDBHandler

func init() {
	defaultMongoDBHandler = newMongoDBHandler()
	defaultMongoDBHandler.setup(config.GlobalServerConf.MongoDBConfig)
}

func newMongoDBHandler() *MongoDBHandler {
	return &MongoDBHandler{
		services: make(map[string]*MongoDBService, 0),
	}
}

func (m *MongoDBHandler) setup(conf map[string]*config.MongoDBConfig) {
	for k, v := range conf {
		newMongoDBService := newMongoDBService()
		if err := newMongoDBService.loadConfig(v); err != nil {
			log.Debug("load mongodb config with key ", k, "with err = ", err)
			continue
		}
		if err := newMongoDBService.openDB(); err != nil {
			log.Debug("open mongodb with key ", k, "with err = ", err)
			continue
		}
		m.services[k] = newMongoDBService
	}
}

func GetCollection(key string) (*mongo.Collection, error) {
	service, ok := defaultMongoDBHandler.services[key]
	if !ok {
		return nil, errors.New("can't find collection with key :" + key + ", make sure db open success")
	}
	return service.collection, nil
}
