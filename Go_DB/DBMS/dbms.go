package DBMS

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	TOMBSTONE    = "`"
	DELIM        = byte(',')
	TOMBSTONELEN = 1
)

func Init() ([]byte, error) {
	dbPath, err := filepath.Abs("Go_DB/Database/db")
	//fmt.Println(dbPath)
	fmt.Println(DELIM)
	db, err := os.ReadFile(dbPath)
	hashTable := InitHashTable(db)
	fmt.Println(LookupHashTable(db, hashTable, 0))
	fmt.Println(LookupHashTable(db, hashTable, 2))
	db, _ = setKey(db, hashTable, 2, "abcdefg")
	db, _ = setKey(db, hashTable, 3, "alilul")
	db, err = deleteKey(db, &hashTable, 2)
	db, err = deleteKey(db, &hashTable, 1)
	_ = saveDb(db)
	return db, err
}

func InitHashTable(db []byte) map[int]int {
	table := map[int]int{}

	i := 0
	n := len(db)
	for i < n {
		offset := i
		key := ""
		for j := i; db[j] != DELIM; j++ {
			key += string(db[j])
		}
		i += len(key) + 1
		keyInt, _ := strconv.Atoi(key)

		valueLen := ""
		for j := i; db[j] != DELIM; j++ {
			valueLen += string(db[j])
		}
		i += len(valueLen) + 1
		valueLenInt, _ := strconv.Atoi(valueLen)

		val := ""
		for j := 0; j < valueLenInt; j++ {
			val += string(db[i+j])
		}

		i += valueLenInt

		if val != TOMBSTONE {
			table[keyInt] = offset
		} else {
			delete(table, keyInt)
		}
	}

	return table
}

func LookupHashTable(db []byte, table map[int]int, key int) (string, error) {
	if offset, ok := table[key]; ok {
		i := offset
		for db[i] != DELIM {
			i++
		}
		i++
		valueLen := ""
		for ; db[i] != DELIM; i++ {
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
		return "", errors.New("key not found")
	}
}

func setKey(db []byte, table map[int]int, key int, value string) ([]byte, error) {
	offset := len(db)
	valueLen := len(value)
	record := fmt.Sprintf("%v,%v,%v", key, valueLen, value)
	db = append(db, []byte(record)...)
	table[key] = offset
	return db, nil
}

func saveDb(db []byte) error {
	dbPath, _ := filepath.Abs("Go_DB/Database/db")
	err := os.WriteFile(dbPath, db, 'w')
	return err
}

func deleteKey(db []byte, table *map[int]int, key int) ([]byte, error) {
	if _, ok := (*table)[key]; ok {
		record := fmt.Sprintf("%v,%v,%s", key, TOMBSTONELEN, TOMBSTONE)
		db = append(db, []byte(record)...)
		delete(*table, key)
		return db, nil
	} else {
		return db, errors.New("key does not exist")
	}
}
