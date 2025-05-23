package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/clarovideo-argentina/goff/interface/redis"
)

func (s *DataStore) Load(dbname string, tblname string) {
	s.Items = make(map[string]StoreItem)

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
					s.Items[name] = StoreItem{
						Value: val,
					}
				}
			}
		}
	}
}

func (s *DataStore) Exec(c redis.Connection, cmdLine [][]byte) redis.Reply {
	for _, v := range cmdLine {
		fmt.Println(string(v[:]))
	}
	cmd := strings.ToUpper(string(cmdLine[0][:]))
	var r redis.Reply
	return r
}
