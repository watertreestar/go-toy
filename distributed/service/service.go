package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Start 启动服务
func Start(ctx context.Context, serviceName, host, port string, registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop.", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
