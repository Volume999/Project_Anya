package main

import (
	"Project_Anya/GoDB/Client"
	"Project_Anya/GoDB/DBMS"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DBABSPATH = "GoDB/Database/db"
)

func main() {
	dbPath, _ := filepath.Abs(DBABSPATH)
	dbms, err := DBMS.Init(dbPath)
	if err != nil {
		panic("Could not initialize database")
	}
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	client := Client.Init(&dbms, reader, writer)
	defer func() {
		_ = dbms.Save()
		_ = writer.Flush()
	}()
	err = client.Run()
	if err != nil {
		fmt.Printf("Writer error: %v\n", err)
	}
}
