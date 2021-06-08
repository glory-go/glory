package mongodb

import (
	"context"
	"fmt"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBService struct {
	conf       *config.MongoDBConfig
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func newMongoDBService() *MongoDBService {
	return &MongoDBService{}
}

func (ms *MongoDBService) loadConfig(conf *config.MongoDBConfig) error {
	ms.conf = conf
	return nil
}

func (ms *MongoDBService) openDB() error {
	DBADDRESS := ms.conf.Host
	DBPORT := ms.conf.Port

	DBUSERNAME := ms.conf.Username
	DBPASSWORD := ms.conf.Password
	var err error
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s", DBUSERNAME, DBPASSWORD, DBADDRESS, DBPORT)
	ms.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		log.Debug("connect mongo with error:", err)
		return err
	}
	ms.db = ms.client.Database(ms.conf.DBName)
	ms.collection = ms.db.Collection(ms.conf.CollectionName)
	return nil
}
