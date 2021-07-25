package minidb

import (
	"errors"
	"github.com/watertreestar/go-toy/minidb/internal"
	"io"
	"os"
	"sync"
)

const (
	DataFileName             = "mini-db.data"
	MergeTempFileName        = "mini-db-merge.data"
	PUT               uint16 = 0
	DEL               uint16 = 1
)

type DBEngine struct {
	indexes map[string]int64 // 内存中的索引信息
	dbFile  *db.File         // 数据文件
	dirPath string           // 数据目录
	mu      sync.RWMutex
}

// NewEngine 创建一个DBEngine
// @param dirPath  数据目录
func NewEngine(dirPath string) (*DBEngine, error) {
	// Create path if not exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	// Create db file by specific data file
	df, err := db.NewDBFile(dirPath, DataFileName)
	if err != nil {
		return nil, err
	}

	engine := &DBEngine{
		dbFile:  df,
		dirPath: dirPath,
		indexes: make(map[string]int64),
	}
	engine.loadIndex()
	return engine, nil
}

// Put 写入数据
func (engine *DBEngine) Put(key []byte, value []byte) (err error) {
	if len(key) == 0 {
		return
	}

	engine.mu.Lock()
	defer engine.mu.Unlock()

	offset := engine.dbFile.Offset
	// 封装成 Entry
	entry := db.NewEntry(key, value, PUT)
	// 追加到数据文件当中
	err = engine.dbFile.AppendEntry(entry)

	// 写到内存
	engine.indexes[string(key)] = offset
	return
}

// Get 取出数据
func (engine *DBEngine) Get(key []byte) (val []byte, err error) {
	if len(key) == 0 {
		return
	}

	engine.mu.RLock()
	defer engine.mu.RUnlock()

	// 从内存当中取出索引信息
	offset, ok := engine.indexes[string(key)]
	// key 不存在
	if !ok {
		return
	}

	// 从磁盘中读取数据
	var e *db.Entry
	e, err = engine.dbFile.ReadEntry(offset)
	if err != nil && err != io.EOF {
		return
	}
	if e != nil {
		val = e.Value
	}
	return
}

// Del 删除数据
func (engine *DBEngine) Del(key []byte) (err error) {
	if len(key) == 0 {
		return
	}

	engine.mu.Lock()
	defer engine.mu.Unlock()
	// 从内存当中取出索引信息
	_, ok := engine.indexes[string(key)]
	// key 不存在，忽略
	if !ok {
		return
	}

	// 封装成 Entry 并写入
	e := db.NewEntry(key, nil, DEL)
	err = engine.dbFile.AppendEntry(e)
	if err != nil {
		return
	}

	// 删除内存中的 key
	delete(engine.indexes, string(key))
	return
}

// loadIndex Load db index from data file in disk
func (engine *DBEngine) loadIndex() error {
	if engine.dbFile == nil {
		return errors.New("DB is nil.Maybe not create")
	}

	var offset int64

	for {
		// Read an entry from data file
		entry, err := engine.dbFile.ReadEntry(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.New("read entry from db file error")
		}

		if entry.Mark == PUT {
			// 设置索引状态
			engine.indexes[string(entry.Key)] = offset
		}
		offset += entry.GetSize()
	}
	return nil
}

// Compress 压缩数据文件
func (engine *DBEngine) Compress() error {
	// Do nothing if no data exist
	if engine.dbFile.Offset == 0 {
		return nil
	}

	var (
		validEntry []*db.Entry
		offset     int64
	)

	for {
		e, err := engine.dbFile.ReadEntry(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// If current offset equals offset in index,it is valid
		if off, ok := engine.indexes[string(e.Key)]; ok && off == offset {
			validEntry = append(validEntry, e)
		}
		offset += e.GetSize()
	}

	if len(validEntry) > 0 {
		// Create a new File instance file to save compressed data
		mdf, err := db.NewDBFile(engine.dirPath, MergeTempFileName)
		if err != nil {
			return err
		}

		defer os.Remove(mdf.File.Name())
		defer mdf.CloseFile()

		// Write entry into new data file
		for _, entry := range validEntry {
			writeOff := mdf.Offset
			err := mdf.AppendEntry(entry)
			if err != nil {
				return err
			}

			// Update index in mem
			engine.indexes[string(entry.Key)] = writeOff
		}

		// Remove old data file and rename merge data file
		os.Remove(engine.dbFile.File.Name())
		os.Rename(mdf.File.Name(), engine.dbFile.Path+string(os.PathSeparator)+engine.dbFile.FileName)

		engine.dbFile = mdf
	}

	return nil
}
