package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gperosio/goff/config"
	"github.com/gperosio/goff/lib/logger"
	RedisServer "github.com/gperosio/goff/redis/server"
	"github.com/gperosio/goff/storage"
	"github.com/gperosio/goff/tcp"
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

	if ds, err := storage.CreateDS(config.StorageProps); err == nil {

		storage.CheckSource(ds, config.StorageProps.CheckTime)

		if err := tcp.ListenAndServeWithSignal(&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		}, RedisServer.MakeHandler(ds)); err != nil {
			logger.Error(err)
		}
	} else {
		logger.Fatal(err)
	}
}
