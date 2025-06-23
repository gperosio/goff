package storage

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/clarovideo-argentina/goff/interface/redis"
	"github.com/clarovideo-argentina/goff/redis/protocol"
)

const MAX_LEVEL_RECURSION = 10

func (s *DataStore) Load(dbname string, tblname string) {
	dot := "."
	basepath := dbname + "/" + tblname

	entries, err := os.ReadDir(basepath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		//tomar el nombre del archivo y guardar como clave, decodear el valor
		var name string
		parts := strings.Split(e.Name(), dot)
		nparts := len(parts)
		if parts[nparts-1] == "json" {
			if nparts > 2 {
				name = strings.Join(parts[0:nparts-1], dot)
			} else {
				name = parts[0]
			}
			//leer archivo
			f := basepath + "/" + e.Name()
			var val map[string]any
			src, err := os.ReadFile(f)
			if err == nil {
				fileInfo, _ := os.Stat(f)
				//parsear json
				if err := json.Unmarshal(src, &val); err == nil {
					s.Items[name] = &StoreItem{
						Value:            val,
						ModificationTime: fileInfo.ModTime(),
					}
				}
			}
		}
	}
}

func (s *DataStore) Exec(c redis.Connection, cmdLine [][]byte) redis.Reply {
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
