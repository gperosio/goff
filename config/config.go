package config

import (
	"os"
	"strconv"
	"time"

	"github.com/clarovideo-argentina/goff/lib/utils"
)

// ServerProperties defines global config properties
type ServerProperties struct {
	// for Public configuration
	RunID       string // runID always different at every exec.
	Bind        string
	Port        int
	Dir         string
	MaxClients  int
	ReplTimeout int
	Storage     string
	Database    string
	Table       string
	CheckTime   int
}

type ServerInfo struct {
	StartUpTime time.Time
}

// Properties holds global config properties
var Properties *ServerProperties
var EachTimeServerInfo *ServerInfo

// SetupConfig read config file and store properties into Properties
func SetupConfig() {
	// A few stats we don't want to reset: server startup time, and peak mem.
	EachTimeServerInfo = &ServerInfo{
		StartUpTime: time.Now(),
	}

	// default config
	Properties = &ServerProperties{
		Bind:      "0.0.0.0",
		Port:      6382,
		RunID:     utils.RandString(40),
		Storage:   "file",
		Database:  "data",
		Table:     "features",
		CheckTime: 300,
	}

	if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		Properties.Port = port
	}

	if chtime, err := strconv.Atoi(os.Getenv("CHECK_TIME")); err == nil {
		Properties.CheckTime = chtime
	}

	if stor := os.Getenv("STORAGE"); stor != "" {
		Properties.Storage = stor
	}

	if db := os.Getenv("DB"); db != "" {
		Properties.Database = db
	}

	if tbl := os.Getenv("TBL"); tbl != "" {
		Properties.Table = tbl
	}
}
