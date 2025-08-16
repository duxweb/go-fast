package main

//go:generate go run github.com/duxweb/go-fast/v2/annotation/annotation-gen

import (
	"example/app/system"
	"example/runtime"

	dux "github.com/duxweb/go-fast/v2"
)

func main() {
	app := dux.New()

	app.RegisterAnnotations(runtime.GetAnnotations())

	app.OnBoot(func() error {
		//fmt.Println("Boot hook")
		return nil
	})

	app.RegisterApp(system.App)

	app.Run()
}
