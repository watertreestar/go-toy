package service

import (
	"context"
	"fmt"
	"github.com/watertreestar/go-toy/distributed/regsitry"
	"log"
	"net/http"
)

// Start 启动服务
func Start(ctx context.Context, host, port string, reg regsitry.Registration, registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, string(reg.ServiceName), host, port)

	err := regsitry.RegisterService(reg)
	if err != nil {
		return nil, err
	}
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
		err := regsitry.UnregisterService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println("unregister service error", err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
