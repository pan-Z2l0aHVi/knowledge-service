package tools

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

type Mongo struct {
}

func (e *Mongo) InitDB() {
	cfg, err := ParseConfigure()
	if err != nil {
		panic(err)
	}
	mongoURL := "mongodb://" +
		cfg.Database.User +
		":" +
		cfg.Database.Pwd +
		"@" +
		cfg.Database.Host +
		":" +
		cfg.Database.Port
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	db = client.Database(cfg.Database.Name)
}

func (e *Mongo) GetDB() *mongo.Database {
	return db
}
