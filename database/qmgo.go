package database

import (
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/qiniu/qmgo"
	"github.com/samber/do"
)

type QmgoService struct {
	client *qmgo.Client
	engine *qmgo.Database
}

func (s *QmgoService) Shutdown() error {
	return s.client.Close(global.Ctx)
}

func Qmgo() *qmgo.Database {
	return do.MustInvoke[*QmgoService](global.Injector).engine
}

func QmgoInit() {
	dbConfig := config.Load("database").GetStringMapString("mongoDB")

	var auth = ""
	if dbConfig["username"] != "" && dbConfig["password"] != "" {
		auth = dbConfig["username"] + ":" + dbConfig["password"] + "@"
	}

	client, err := qmgo.NewClient(global.Ctx, &qmgo.Config{Uri: "mongodb://" + auth + dbConfig["host"] + ":" + dbConfig["port"]})
	if err != nil {
		panic("qmgo error :" + err.Error())
	}

	do.ProvideValue[*QmgoService](global.Injector, &QmgoService{
		client: client,
		engine: client.Database(dbConfig["dbname"]),
	})
}
