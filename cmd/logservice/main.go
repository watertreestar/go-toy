package main

import (
	"context"
	"fmt"
	"github.com/watertreestar/go-toy/distributed/log"
	"github.com/watertreestar/go-toy/distributed/regsitry"
	"github.com/watertreestar/go-toy/distributed/service"
	stlog "log"
)

func main() {
	log.Init("./distributed.log")
	host, port := "localhost", "9100"
	r := regsitry.Registration{
		"LogService",
		fmt.Sprintf("http://%s:%s", host, port),
	}
	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service")
}
