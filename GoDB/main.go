package main

import (
	"Project_Anya/GoDB/Client"
	"Project_Anya/GoDB/DBMS"
	"Project_Anya/GoDB/Evaluation"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DBAbsPath = "GoDB/Database/db"
	RunMode   = "Eval" // Run | Eval
)

func run() {
	dbPath, _ := filepath.Abs(DBAbsPath)
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

func main() {
	switch RunMode {
	case "Run":
		run()
	case "Eval":
		Evaluation.RunEval()
	default:
		panic("invalid run mode")
	}

}
