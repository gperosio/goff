package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

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
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	addr := "0.0.0.0:" + getenv("PORT", "6382")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind to %s\n", addr)
		os.Exit(1)
	}

	fmt.Print(banner)
	log.Printf("Starting Goff server at %s\n", addr)

	CreateDS(getenv("STORAGE", "file"), getenv("ST_DB", "data"), getenv("ST_TBL", "features"))

	/*
		go func() {
			for {
				//leer mayor a fecha_ultima sino todo
				// si resultado, hacer cambio
				//actualizar fecha ultima
			}
		}()
	*/
	// Listen for inputs and respond
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("failed to accept conn: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 2014) // store out stuff somewhere
			_ = conn.SetReadDeadline(time.Now().Add(ReadWriteTimeout))

			for {
				len, err := conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Printf("Error reading: %#v\n", err)
					}
					break
				}

				// Parse the command out
				command := buf[:len]
				response := Parse(command)

				// Write the response back to the connection
				_, responseErr := conn.Write([]byte(response + "\r\n"))
				if responseErr != nil {
					fmt.Printf("Error reading: %#v\n", err)
					break
				}
				_ = conn.SetReadDeadline(time.Now().Add(ReadWriteTimeout))
			}
		}(conn)
	}
}
