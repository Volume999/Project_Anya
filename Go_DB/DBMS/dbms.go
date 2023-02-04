package DBMS

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func Init() ([]byte, error) {
	dbPath, err := filepath.Abs("Go_DB/Database/db")
	//fmt.Println(dbPath)
	db, err := os.ReadFile(dbPath)
	hashTable := InitHashTable(db)
	fmt.Println(LookupHashTable(db, hashTable, 0))
	fmt.Println(LookupHashTable(db, hashTable, 2))
	return db, err
}

func InitHashTable(db []byte) map[int]int {
	table := map[int]int{}

	i := 0
	n := len(db)
	for i < n {
		offset := i
		key := ""
		for j := i; db[j] != 44; j++ {
			key += string(db[j])
		}
		i += len(key) + 1
		keyInt, _ := strconv.Atoi(key)

		valueLen := ""
		for j := i; db[j] != 44; j++ {
			valueLen += string(db[j])
		}
		i += len(valueLen) + 1
		valueLenInt, _ := strconv.Atoi(valueLen)

		i += valueLenInt

		table[keyInt] = offset
	}

	return table
}

func LookupHashTable(db []byte, table map[int]int, key int) (string, error) {
	if offset, ok := table[key]; ok {
		i := offset
		for db[i] != 44 {
			i++
		}
		i++
		valueLen := ""
		for ; db[i] != 44; i++ {
			valueLen += string(db[i])
		}
		valueLenInt, _ := strconv.Atoi(valueLen)
		i++

		val := ""
		for j := 0; j < valueLenInt; j++ {
			val += string(db[i+j])
		}

		return val, nil
	} else {
		return "", errors.New("Key not found")
	}
}
