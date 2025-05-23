package storage

import (
	"sync"
	"time"
)

type StoreItem struct {
	Value            map[string]any
	ModificationTime time.Time
	Mutex            sync.RWMutex
}

type DataStore struct {
	Items map[string]StoreItem
}

func CreateDS(storage_type string, dbname string, tblname string) *DataStore {
	ds := DataStore{Items: make(map[string]StoreItem)}
	ds.Load(dbname, tblname)
	return &ds
}
