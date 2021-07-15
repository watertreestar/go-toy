package main

import (
	"context"
	"fmt"
	"github.com/watertreestar/go-toy/distributed/log"
	"github.com/watertreestar/go-toy/distributed/service"
	stlog "log"
)

func main() {
	log.Init("./distributed.log")
	host, port := "localhost", "9100"
	ctx, err := service.Start(context.Background(),
		"LogService",
		host,
		port,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service")
}
