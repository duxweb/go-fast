package database

import (
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/gookit/event"
	"github.com/qiniu/qmgo"
	"github.com/samber/do"
)

func QmgoInit() {
	dbConfig := config.Get("database").GetStringMapString("mongoDB")

	var auth = ""
	if dbConfig["username"] != "" && dbConfig["password"] != "" {
		auth = dbConfig["username"] + ":" + dbConfig["password"] + "@"
	}

	client, err := qmgo.NewClient(global.Ctx, &qmgo.Config{Uri: "mongodb://" + auth + dbConfig["host"] + ":" + dbConfig["port"]})
	if err != nil {
		panic("qmgo error :" + err.Error())
	}

	do.ProvideValue[*qmgo.Database](nil, client.Database(dbConfig["dbname"]))

	event.On("app.close", event.ListenerFunc(func(e event.Event) error {
		return client.Close(global.Ctx)
	}))
}
