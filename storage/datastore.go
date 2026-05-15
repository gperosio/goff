package storage

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gperosio/goff/config"
	"github.com/gperosio/goff/interface/redis"
	"github.com/gperosio/goff/redis/protocol"
)

const MAX_LEVEL_RECURSION = 10

type StoreItem struct {
	Value            map[string]any
	ModificationTime time.Time
	Mutex            sync.RWMutex
}

type DataStore struct {
	Items        map[string]*StoreItem
	Source       DataproviderSource
	MaxModifTime time.Time
}

type DataproviderSource interface {
	Init() error
	Load(items map[string]*StoreItem, maxModifTime time.Time) (time.Time, error)
	Close()
}

var store = make(map[string]*StoreItem)

// crea datasource
func CreateDS(sp *config.StorageProperties) (*DataStore, error) {
	var ds *DataStore
	var maxtime time.Time
	var err error

	if sp.Type == "file" {
		ds = &DataStore{Items: store, Source: FileDataSource{LocalDir: sp.LocalDir}}
	}
	if sp.Type == "postgres" {
		ds = &DataStore{Items: store, Source: PgDataSource{Host: sp.Host, Port: sp.Port, Database: sp.Database, User: sp.User, Passwd: sp.Passwd, Table: sp.Table}}
	}
	if err := ds.Source.Init(); err != nil {
		return nil, err
	}
	if maxtime, err = ds.Source.Load(ds.Items, ds.MaxModifTime); err != nil {
		return nil, err
	}
	ds.MaxModifTime = maxtime
	return ds, nil
}

func DatastoreSynchronizer(ds *DataStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				ds.ReLoad()
			}
		}
	}()
}

func (s *DataStore) ReLoad() {
	s.MaxModifTime, _ = s.Source.Load(s.Items, s.MaxModifTime)
}

func (s *DataStore) Exec(client redis.Connection, cmdLine [][]byte) redis.Reply {
	cmd := strings.ToUpper(string(cmdLine[0][:]))
	switch cmd {
	case "PING":
		return protocol.MakeBulkReply([]byte("PONG"))

	case "GET":
		if len(cmdLine) < 2 {
			return protocol.MakeErrReply("Missing key")
		}
		item := s.Items[string(cmdLine[1][:])]
		if item == nil {
			return &protocol.NullBulkReply{}
		}
		item.Mutex.RLock()
		defer item.Mutex.RUnlock()
		msg, _ := json.Marshal(item.Value)
		return protocol.MakeBulkReply(msg)

	case "HGET":
		if len(cmdLine) < 3 {
			return protocol.MakeErrReply("Missing key")
		}
		item := s.Items[string(cmdLine[1][:])]
		if item == nil {
			return &protocol.NullBulkReply{}
		}
		item.Mutex.RLock()
		defer item.Mutex.RUnlock()

		parts := strings.Split(string(cmdLine[2]), ".")
		if len(parts) > MAX_LEVEL_RECURSION {
			return protocol.MakeErrReply("Depth level is higher, max " + strconv.Itoa(MAX_LEVEL_RECURSION))
		}

		// val := findKey(item, parts)
		val := recurseMap(item.Value, parts)
		if val == nil {
			return &protocol.NullBulkReply{}
		}
		msg, _ := json.Marshal(val)
		return protocol.MakeBulkReply(msg)

	case "KEYS":
		var out [][]byte
		for k := range s.Items {
			out = append(out, []byte(k))
		}
		return protocol.MakeMultiBulkReply(out)

	default:
		return protocol.MakeErrReply("ERR unknown command '" + cmd + "'")
	}
}

func unshift(alist *[]string) string {
	rv := (*alist)[0]
	*alist = (*alist)[1:]
	return rv
}

func recurseMap(data any, keys []string) any {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if len(keys) == 0 {
		return data
	}

	if val, ok := data.(map[string]any)[keys[0]]; ok {
		unshift(&keys)
		return recurseMap(val, keys)
	}
	return nil
}
