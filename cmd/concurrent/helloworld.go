package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)
	for i := 0; i < 5000; i++ {
		go printHello(i, ch)
	}

	for {
		msg := <-ch
		fmt.Println(msg)
	}
	// time.Sleep(time.Millisecond * 200)
}

func printHello(i int, ch chan string) {
	for {
		ch <- fmt.Sprintf("Hello world from %d\n", i)
	}
}
