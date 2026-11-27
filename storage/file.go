package storage

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type FileDataSource struct {
	LocalDir string
}

func (f FileDataSource) Init() error {
	//verificar que el path existe
	if _, err := os.Stat(f.LocalDir); os.IsNotExist(err) {
		return err
	}
	return nil
}

func (f FileDataSource) Load(items map[string]*StoreItem, maxModifTime time.Time) (time.Time, error) {
	var maxtime time.Time
	dot := "."

	entries, err := os.ReadDir(f.LocalDir)
	if err != nil {
		return maxtime, err
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
			f := f.LocalDir + "/" + e.Name()
			var val map[string]any
			src, err := os.ReadFile(f)
			if err == nil {
				fileInfo, _ := os.Stat(f)
				modtime := fileInfo.ModTime()
				if items[name] == nil || items[name] != nil && modtime.Unix() > maxModifTime.Unix() {
					//parsear json
					if err := json.Unmarshal(src, &val); err == nil {
						items[name] = &StoreItem{
							Value:            val,
							ModificationTime: modtime,
						}
					}
					if modtime.Unix() > maxtime.Unix() {
						maxtime = modtime
					}
				}
			}
		}
	}
	return maxtime, nil
}

func (f FileDataSource) Close() {}
