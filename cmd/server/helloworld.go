package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "<p>hello,world!</p>")
	})

	http.ListenAndServe(":8888", nil)
}
