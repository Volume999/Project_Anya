package Tests

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	EmptyDb    = "Mock_DBs/empty_db_test"
	NonEmptyDb = "Mock_DBs/nonempty_db_test"
	DbForSave  = "Mock_DBs/db_for_save_test"
)

var EmptyDBAbsPath string
var NonEmptyDBAbsPath string
var DbForSaveAbsPath string

func TestMain(m *testing.M) {
	path, err := filepath.Abs(EmptyDb)
	if err != nil {
		panic("empty db path is invalid")
	}
	EmptyDBAbsPath = path
	path, err = filepath.Abs(NonEmptyDb)
	if err != nil {
		panic("nonempty db path is invalid")
	}
	NonEmptyDBAbsPath = path
	path, err = filepath.Abs(DbForSave)
	if err != nil {
		panic("db for save path is invalid")
	}
	DbForSaveAbsPath = path
	runTests := m.Run()
	os.Exit(runTests)
}
