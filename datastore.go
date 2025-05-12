package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type StoreItem struct {
	Value            map[string]any
	ModificationTime time.Time
	Mutex            sync.Mutex
}

var store = make(map[string]*StoreItem)

func CreateDS(storage_type string, dbname string, tblname string) {
	if storage_type == "file" {
		LoadFile(dbname, tblname)
	}
}

func findKey(item *StoreItem, keys []string) any {
	// Item exists - enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	if len(keys) == 0 {
		return item.Value
	}

	//fijo por niveles
	if len(keys) == 1 {
		if val, ok := item.Value[keys[0]]; ok {
			return val
		}
		return nil
	}

	if len(keys) == 2 {
		if val, ok := item.Value[keys[0]].(map[string]any)[keys[1]]; ok {
			return val
		}
		return nil
	}

	if len(keys) == 3 {
		if val, ok := item.Value[keys[0]].(map[string]any)[keys[1]].(map[string]any)[keys[2]]; ok {
			return val
		}
		return nil
	}

	return nil
}

func LoadFile(dbname string, tblname string) {
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
				//parsear json
				if err := json.Unmarshal(src, &val); err == nil {
					store[name] = &StoreItem{
						Value: val,
					}
				}
			}
		}
	}
}
