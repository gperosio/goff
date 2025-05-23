package main

import (
	"fmt"
	"os"
	"time"

	"github.com/clarovideo-argentina/goff/config"
	"github.com/clarovideo-argentina/goff/lib/logger"
	RedisServer "github.com/clarovideo-argentina/goff/redis/server"
	"github.com/clarovideo-argentina/goff/storage"
	"github.com/clarovideo-argentina/goff/tcp"
	"github.com/joho/godotenv"
)

const MaxLineLenBytes = 2048
const ReadWriteTimeout = time.Minute

var banner = `
                ________
   ____ _____  / __/ __/
  / __ ` + "`" + `/ __ \/ /_/ /_  
 / /_/ / /_/ / __/ __/  
 \__, /\____/_/ /_/     
/____/
         
`

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	config.SetupConfig()

	fmt.Print(banner)
	//logger.Info(fmt.Sprintf("Starting Goff server at %s:%d", config.Properties.Bind, config.Properties.Port))

	//hacer que devuelva un puntero para pasarselo a la gorutina de revision y al servicio TCP
	ds := storage.CreateDS(config.Properties.Storage, config.Properties.Database, config.Properties.Table)
	/*
		go func() {
			for {
				//leer mayor a fecha_ultima sino todo
				// si resultado, hacer cambio
				//actualizar fecha ultima
			}
		}()
	*/
	if err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, RedisServer.MakeHandler(ds)); err != nil {
		logger.Error(err)
	}
}
