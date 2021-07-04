package main

import (
	"bufio"
	"fmt"
	"github.com/watertreestar/go-toy/sort/pipeline"
	"os"
	"strconv"
)

func main() {
	// p := createPipeline("big.in", 800000000, 8)
	p := createNetworkPipeline("big.in", 800000000, 8)

	writeFile(p, "big.out")
	printFile("big.out")

	// Test results
	// 4 chunkCount : 28s
	// 8 chunkCount : 26s
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	p := pipeline.ReaderSource(file, -1)
	count := 0
	for v := range p {
		if count > 100 {
			break
		}
		fmt.Println(v)
		count++
	}
}

// Write data to writer
func writeFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	pipeline.WriterSink(writer, p)
}

func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount

	var sortResults []<-chan int
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(chunkSize*i), 0)

		pipeline.InitTime()
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		result := pipeline.MemSort(source)
		sortResults = append(sortResults, result)
	}
	return pipeline.MergeN(sortResults...)
}

// Sink data via network
func createNetworkPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount

	var sortAddr []string
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(chunkSize*i), 0)

		pipeline.InitTime()
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		result := pipeline.MemSort(source)

		addr := ":" + strconv.Itoa(9000+i)
		pipeline.NetworkSink(addr, result)
		// Send to merge node by network
		sortAddr = append(sortAddr, addr)
	}

	// Merge node to connect sort node
	var sortResults []<-chan int
	for _, sortAddr := range sortAddr {
		sortResults = append(sortResults, pipeline.NetworkSource(sortAddr))
	}
	return pipeline.MergeN(sortResults...)
}
