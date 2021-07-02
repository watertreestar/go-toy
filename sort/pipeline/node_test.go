package pipeline

import (
	"os"
	"testing"
)

func TestArraySource(t *testing.T) {
	p := ArraySource(3, 2, 6, 7, 4)
	for {
		if num, ok := <-p; ok {
			t.Log(num)
		} else {
			break
		}
	}
}

func TestMemSort(t *testing.T) {
	p := MemSort(ArraySource(3, 2, 6, 7, 4))
	for {
		if num, ok := <-p; ok {
			t.Log(num)
		} else {
			break
		}
	}
}

func TestMerge(t *testing.T) {
	p := Merge(MemSort(ArraySource(3, 2, 6, 7, 4)),
		MemSort(ArraySource(23, 5, 7, 12, 0, 6, 9, 7, 3, 80, 26)))

	for num := range p {
		t.Log(num)
	}
}

func TestReaderSource(t *testing.T) {
	t.Run("Small Data", func(t *testing.T) {
		file, err := os.Create("small.in")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		p := RandomSource(200)
		WriterSink(file, p)

		file, err = os.Open("small.in")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		p = ReaderSource(file, -1)
		for v := range p {
			t.Log(v)
		}
	})

	t.Run("Big Data", func(t *testing.T) {
		const filename = "big.in"
		const n = 800000
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		p := RandomSource(n)
		WriterSink(file, p)

		file, err = os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		p = ReaderSource(file, -1)
		//for v := range p {
		//	t.Log(v)
		//}
	})

}
