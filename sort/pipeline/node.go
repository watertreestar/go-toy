package pipeline

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var startTime time.Time

func InitTime() {
	startTime = time.Now()
}

func ArraySource(a ...int) chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

func MemSort(in <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		var a []int
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read done:", time.Now().Sub(startTime))
		sort.Ints(a)
		fmt.Println("Mem sort done:", time.Now().Sub(startTime))
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// Merge 归并排好序的两路数据
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge done:", time.Now().Sub(startTime))
	}()

	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[0:m]...), MergeN(inputs[m:]...))
}

// ReaderSource 从Reader来的数据
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		bytesReader := 0
		bufferReader := bufio.NewReader(reader)
		for {
			n, err := bufferReader.Read(buffer)
			bytesReader += n
			if n > 0 {
				v := binary.BigEndian.Uint64(buffer)
				out <- int(v)
			}
			if err != nil || (chunkSize != -1 && bytesReader >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	bufferWriter := bufio.NewWriter(writer)
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		bufferWriter.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}
