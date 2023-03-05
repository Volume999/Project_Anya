package Tests

import (
	"Project_Anya/GoDB/DBMS"
	"testing"
)

func TestInitEmptyDB(t *testing.T) {
	dbms, err := DBMS.Init(EmptyDBAbsPath)
	if err != nil {
		t.Error(err)
	} else {
		if dbms.Size() != 0 {
			t.Errorf("Empty DB Size is %v instead of %v", dbms.Size(), 0)
		}
	}
}

// TODO: Expand this test, and move it after saving DB
func TestInitNonEmptyDB(t *testing.T) {
	dbms, err := DBMS.Init(NonEmptyDBAbsPath)
	if err != nil {
		t.Error(err)
	} else {
		if dbms.Size() != 2 {
			t.Errorf("DB Size is %v instead of %v", dbms.Size(), 2)
		}
	}
}
