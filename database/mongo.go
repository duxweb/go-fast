package database

import (
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	client *mongo.Client
	engine *mongo.Database
}

func (s *MongoService) Shutdown() error {
	return s.client.Disconnect(global.CtxBackground)
}

func Mongo(name ...string) *mongo.Database {
	n := "default"
	if len(name) > 0 {
		n = name[0]
	}

	client, err := do.InvokeNamed[*MongoService](global.Injector, "mongodb."+n)
	if err != nil {
		client = NewMongo(n)
		do.ProvideNamedValue[*MongoService](global.Injector, "mongodb."+n, client)
	}
	return client.engine
}

func NewMongo(name string) *MongoService {
	dbConfig := config.Load("database").GetStringMapString("mongodb.drivers." + name)

	var auth = ""
	if dbConfig["username"] != "" && dbConfig["password"] != "" {
		auth = dbConfig["username"] + ":" + dbConfig["password"] + "@"
	}
	client, err := mongo.Connect(global.CtxBackground, options.Client().ApplyURI("mongodb://"+auth+dbConfig["host"]+":"+dbConfig["port"]))
	if err != nil {
		panic("qmgo error :" + err.Error())
	}

	return &MongoService{
		client: client,
		engine: client.Database(dbConfig["dbname"]),
	}
}
