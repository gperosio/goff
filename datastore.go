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
	Mutex            sync.RWMutex
}

var store = make(map[string]*StoreItem)

func CreateDS(storage_type string, dbname string, tblname string) {
	if storage_type == "file" {
		LoadFile(dbname, tblname)
	}
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
