package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDBFile(t *testing.T) {
	_, err := NewDBFile("./", "db.data")
	if err != nil {
		t.Error("error new db file")
	}
}

func TestFile_AppendEntry(t *testing.T) {
	df, err := NewDBFile("./", "db.data")
	if err != nil {
		t.Error("error new db file")
	}

	entry := NewEntry([]byte(`name`), []byte("young"), 10)
	df.AppendEntry(entry)

	e, err := df.ReadEntry(0)
	if err != nil {
		t.Error("read entry from file error")
	}
	assert.Equal(t, []byte("young"), e.Value)
	assert.Equal(t, []byte("name"), e.Key)
}

func BenchmarkFile_AppendEntry(b *testing.B) {
	df, err := NewDBFile("./", "db.data")
	if err != nil {
		b.Error("error new db file")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := NewEntry([]byte(`name`), []byte("young"), 10)
		df.AppendEntry(entry)
	}

}

func BenchmarkTestFile_ReadEntry(b *testing.B) {
	df, err := NewDBFile("./", "db.data")
	if err != nil {
		b.Error("error new db file")
	}

	for i := 0; i < b.N; i++ {
		entry := NewEntry([]byte(`name`), []byte("young"), 10)
		df.AppendEntry(entry)
	}
	b.ResetTimer()
	var offset int64 = 0
	for i := 0; i < b.N; i++ {
		e, err := df.ReadEntry(offset)
		if err != nil {
			panic(err)
		}
		assert.Equal(b, []byte("young"), e.Value)
		assert.Equal(b, []byte("name"), e.Key)
		offset += e.GetSize()
	}

}
