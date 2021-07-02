package service

import (
	"errors"
	db "github.com/watertreestar/go-toy/minidb/internal"
	"io"
	"os"
	"sync"
)

const (
	DataFileName        = "mini-db.data"
	PUT          uint16 = iota
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
	// todo Load index from disk

	return engine, nil
}

// loadIndex Load db index from data file in disk
func (engine *DBEngine) loadIndex() error {
	if engine.dbFile == nil {
		return errors.New("DB is nil.Maybe not create")
	}

	var offset int64

	for {
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
