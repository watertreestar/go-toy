package main

import (
	"context"
	"fmt"
	"github.com/watertreestar/go-toy/distributed/regsitry"
	"log"
	"net/http"
)

func main() {
	http.Handle("/services", &regsitry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = regsitry.ServerPort

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("Register service serve error: %s", err)
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop.", "Registry Service")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down registry service")
}
