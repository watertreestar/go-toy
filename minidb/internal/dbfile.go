package db

import (
	"os"
)

type File struct {
	Path     string
	File     *os.File
	Offset   int64
	FileName string
}

// AppendEntry  追加Entry到文件
func (dbFile *File) AppendEntry(entry *Entry) error {
	enc, err := entry.Encode()
	if err != nil {
		return err
	}
	_, err = dbFile.File.Write(enc)
	// Increase current offset to write data
	dbFile.Offset += entry.GetSize()
	return nil
}

// CloseFile close data file that File instance holds on
func (dbFile *File) CloseFile() {
	dbFile.File.Close()
}

// ReadEntry 从file中指定的offset处开始读取
func (dbFile *File) ReadEntry(offset int64) (*Entry, error) {
	// 1. Reader entry header first
	buf := make([]byte, EntryHeaderSize)
	if _, err := dbFile.File.ReadAt(buf, offset); err != nil {
		return nil, err
	}
	e, err := Decode(buf)
	if err != nil {
		return nil, err
	}

	offset += int64(EntryHeaderSize)

	// Extract key size and value size from entry header
	// 2. Read key from entry file
	if e.KeySize > 0 {
		key := make([]byte, e.KeySize)
		if _, err = dbFile.File.ReadAt(key, offset); err != nil {
			return nil, err
		}
		e.Key = key
	}

	offset += int64(e.KeySize)
	// 3. Read value from entry file
	if e.ValueSize > 0 {
		value := make([]byte, e.ValueSize)
		if _, err = dbFile.File.ReadAt(value, offset); err != nil {
			return nil, err
		}
		e.Value = value
	}

	return e, nil
}

// Merge 创建一个merge数据文件
func (dbFile *File) Merge(mergeFileName string) error {
	return nil
}

// NewDBFile 创建一个DBFile
func NewDBFile(path, filename string) (*File, error) {
	finalPath := path + string(os.PathSeparator) + filename

	file, err := os.OpenFile(finalPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(finalPath)
	if err != nil {
		return nil, err
	}
	// Offset is zero when file create
	return &File{Offset: stat.Size(), File: file, Path: path, FileName: filename}, nil
}
