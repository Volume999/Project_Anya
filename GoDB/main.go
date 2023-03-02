package main

import (
	"Project_Anya/GoDB/Client"
	"Project_Anya/GoDB/DBMS"
	"bufio"
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
	client := Client.Init(&dbms, reader)
	defer func() {
		_ = dbms.Save()
	}()
	client.Run()
}
