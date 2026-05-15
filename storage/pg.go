package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type PgDataSource struct {
	Host     string
	Port     int
	Database string
	User     string
	Passwd   string
	Table    string
	Conn     *pgx.Conn
}

func (p PgDataSource) Init() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.User, p.Passwd, p.Host, p.Port, p.Database)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	p.Conn = conn
	return nil
}

func (p PgDataSource) Load(items map[string]*StoreItem, maxModifTime time.Time) (time.Time, error) {
	var maxtime time.Time
	var err error
	var q string

	if maxModifTime.IsZero() {
		//si es la primera carga, traigo todo
		q = fmt.Sprintf("SELECT id, document, date_loaded FROM %s order by date_loaded", p.Table)
	} else {
		q = fmt.Sprintf("SELECT id, document, date_loaded FROM %s WHERE date_loaded > $1 order by date_loaded", p.Table)
	}
	rows, err := p.Conn.Query(context.Background(), q, maxModifTime)
	if err != nil {
		log.Fatalf("Database query failed because %s", err)
	}
	defer rows.Close()

	docInfo := make(map[string]time.Time)
	for rows.Next() {
		var doc_id string
		var date_loaded time.Time
		err = rows.Scan(&doc_id, &date_loaded)
		if err != nil {
			log.Fatalf("Failed to retrieve row because %s", err)
		}
		maxtime = date_loaded
		docInfo[doc_id] = date_loaded
	}

	//traigo de a uno
	for key, date_loaded := range docInfo {
		var src []byte
		var val map[string]any
		row := p.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT document FROM %s WHERE id = $1", p.Table), key)
		err = row.Scan(&src)
		if err != nil {
			log.Fatalf("Failed to retrieve document for ID %s because %s", key, err)
		}
		if err := json.Unmarshal(src, &val); err == nil {
			items[key] = &StoreItem{
				Value:            val,
				ModificationTime: date_loaded,
			}
		}
	}
	return maxtime, err
}

func (p PgDataSource) Close() {
	p.Conn.Close(context.Background())
}
