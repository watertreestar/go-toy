package db

import (
	"encoding/binary"
)

const (
	EntryHeaderSize uint32 = 10
)

type Entry struct {
	Key       []byte
	Value     []byte
	KeySize   uint32
	ValueSize uint32
	Mark      uint16
}

// NewEntry 创建一个Entry
func NewEntry(key, value []byte, mark uint16) *Entry {
	return &Entry{
		Key:       key,
		Value:     value,
		Mark:      mark,
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
	}
}

// GetSize 获取Entry的大小
func (e *Entry) GetSize() int64 {
	return int64(EntryHeaderSize + e.KeySize + e.ValueSize)
}

// Encode 编码 Entry，返回字节数组
func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, e.GetSize())
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	binary.BigEndian.PutUint16(buf[8:10], e.Mark)
	copy(buf[EntryHeaderSize:EntryHeaderSize+e.KeySize], e.Key)
	copy(buf[EntryHeaderSize+e.KeySize:], e.Value)
	return buf, nil
}

// Decode 解码 buf 字节数组，返回 Entry
func Decode(buf []byte) (*Entry, error) {
	ks := binary.BigEndian.Uint32(buf[0:4])
	vs := binary.BigEndian.Uint32(buf[4:8])
	mark := binary.BigEndian.Uint16(buf[8:10])
	return &Entry{KeySize: ks, ValueSize: vs, Mark: mark}, nil
}
