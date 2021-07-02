package main

import (
	"bufio"
	"fmt"
	"github.com/watertreestar/go-toy/sort/pipeline"
	"os"
)

func main() {
	p := createPipeline("big.in", 800000, 8)
	writeFile(p, "big.out")
	printFile("big.out")
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	p := pipeline.ReaderSource(file, -1)

	for v := range p {
		fmt.Println(v)
	}
}

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

		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		result := pipeline.MemSort(source)
		sortResults = append(sortResults, result)
	}
	return pipeline.MergeN(sortResults...)
}
