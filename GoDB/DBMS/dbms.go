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

type Dbms struct {
	hashTable map[int]int
	db        []byte
}

func Init(dbPath string) (Dbms, error) {
	db, err := os.ReadFile(dbPath)

	dbms := Dbms{
		hashTable: map[int]int{},
		db:        db,
	}

	dbms.initHashTable()
	return dbms, err
}

func (dbms *Dbms) initHashTable() {
	i := 0
	n := len(dbms.db)
	for i < n {
		offset := i
		key := ""
		for j := i; dbms.db[j] != DELIM; j++ {
			key += string(dbms.db[j])
		}
		i += len(key) + 1
		keyInt, _ := strconv.Atoi(key)

		valueLen := ""
		for j := i; dbms.db[j] != DELIM; j++ {
			valueLen += string(dbms.db[j])
		}
		i += len(valueLen) + 1
		valueLenInt, _ := strconv.Atoi(valueLen)

		val := ""
		for j := 0; j < valueLenInt; j++ {
			val += string(dbms.db[i+j])
		}

		i += valueLenInt

		if val != TOMBSTONE {
			dbms.hashTable[keyInt] = offset
		} else {
			delete(dbms.hashTable, keyInt)
		}
	}
}

func (dbms *Dbms) Get(key int) (string, error) {
	offset, err := dbms.getKeyOffset(key)
	if err == nil {
		val, err := dbms.getValueByOffset(offset)
		return val, err
	} else {
		return "", err
	}
}

func (dbms *Dbms) getKeyOffset(key int) (int, error) {
	if offset, ok := dbms.hashTable[key]; ok {
		return offset, nil
	} else {
		return -1, errors.New("key not found")
	}
}

func (dbms *Dbms) getValueByOffset(offset int) (string, error) {
	i := offset
	for dbms.db[i] != DELIM {
		i++
	}
	i++
	valueLen := ""
	for ; dbms.db[i] != DELIM; i++ {
		valueLen += string(dbms.db[i])
	}
	valueLenInt, _ := strconv.Atoi(valueLen)
	i++

	val := ""
	for j := 0; j < valueLenInt; j++ {
		val += string(dbms.db[i+j])
	}

	return val, nil
}

func (dbms *Dbms) Set(key int, value string) {
	offset := len(dbms.db)
	valueLen := len(value)
	record := fmt.Sprintf("%v,%v,%v", key, valueLen, value)
	dbms.db = append(dbms.db, []byte(record)...)
	dbms.hashTable[key] = offset
}

// TODO: Add Filepath to the property of DB
func (dbms *Dbms) Save() error {
	dbPath, _ := filepath.Abs("GoDB/Database/db")
	err := os.WriteFile(dbPath, dbms.db, 'w')
	return err
}

func (dbms *Dbms) Delete(key int) error {
	if _, ok := (dbms.hashTable)[key]; ok {
		record := fmt.Sprintf("%v,%v,%s", key, TOMBSTONELEN, TOMBSTONE)
		dbms.db = append(dbms.db, []byte(record)...)
		delete(dbms.hashTable, key)
		return nil
	} else {
		return errors.New("key does not exist")
	}
}
