package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gperosio/goff/lib/utils"
)

type StorageProperties struct {
	Type      string
	Host      string
	Port      int
	Database  string
	Table     string
	User      string
	Passwd    string
	LocalDir  string
	CheckTime time.Duration
}

// ServerProperties defines global config properties
type ServerProperties struct {
	// for Public configuration
	RunID       string // runID always different at every exec.
	Bind        string
	Port        int
	Dir         string
	MaxClients  int
	ReplTimeout int
}

type ServerInfo struct {
	StartUpTime time.Time
}

// Properties holds global config properties
var Properties *ServerProperties
var EachTimeServerInfo *ServerInfo
var StorageProps *StorageProperties

// SetupConfig read config file and store properties into Properties
func SetupConfig() {
	// A few stats we don't want to reset: server startup time, and peak mem.
	EachTimeServerInfo = &ServerInfo{
		StartUpTime: time.Now(),
	}

	// default config
	Properties = &ServerProperties{
		Bind:  "0.0.0.0",
		Port:  6382,
		RunID: utils.RandString(40),
	}

	if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil && port > 0 {
		Properties.Port = port
	}

	StorageProps = &StorageProperties{
		Type:      "file",
		LocalDir:  "data/features",
		CheckTime: 300,
	}

	if stor := os.Getenv("STORAGE"); stor != "" {
		StorageProps.Type = stor
	}

	if localdir := os.Getenv("LOCAL_DIR"); localdir != "" {
		StorageProps.LocalDir = localdir
	}

	if db := os.Getenv("DB_NAME"); db != "" {
		StorageProps.Database = db
	}

	if tbl := os.Getenv("DB_TBL"); tbl != "" {
		StorageProps.Table = tbl
	}

	if user := os.Getenv("DB_USER"); user != "" {
		StorageProps.User = user
	}

	if pass := os.Getenv("DB_PASSWD"); pass != "" {
		StorageProps.Passwd = pass
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		StorageProps.Host = host
	}

	if port, err := strconv.Atoi(os.Getenv("DB_PORT")); err == nil && port > 0 {
		StorageProps.Port = port
	}
	if chtime, err := strconv.Atoi(os.Getenv("CHECK_TIME")); err == nil && chtime > 0 {
		StorageProps.CheckTime = time.Duration(chtime) * time.Second
	}
}
