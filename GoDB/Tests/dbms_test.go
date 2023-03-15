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

func TestGet(t *testing.T) {
	dbms, _ := DBMS.Init(NonEmptyDBAbsPath)
	if val, err := dbms.Get(2); err == nil {
		if val != "AliPizza" {
			t.Errorf("Incorrect value, %v instead of %v", val, "AliPizza")
		}
	} else {
		t.Error(err)
	}
}

func TestGetWithMultipleVersions(t *testing.T) {
	dbms, _ := DBMS.Init(NonEmptyDBAbsPath)
	if val, _ := dbms.Get(4); val != "ab" {
		t.Errorf("Incorrect value, %v instead of %v", val, "ab")
	}
}

func TestGetNotFound(t *testing.T) {
	dbms, _ := DBMS.Init(NonEmptyDBAbsPath)
	if val, err := dbms.Get(1); err == nil || val != "" {
		t.Errorf("Error not caught, err: %v, val: %v", err, val)
	} else {
		if err.Error() != "key not found" {
			t.Errorf("Incorrect error, %v instead of %v", err.Error(), "key not found")
		}
	}
}

func TestSet(t *testing.T) {
	dbms, _ := DBMS.Init(EmptyDBAbsPath)
	dbms.Set(1, "testing")
	if val, err := dbms.Get(1); err != nil || val != "testing" {
		t.Errorf("Incorrect value returned after set, Error: %v, Value: %v", err, val)
	}
}

func TestSetMultiple(t *testing.T) {
	dbms, _ := DBMS.Init(EmptyDBAbsPath)
	dbms.Set(1, "testing")
	dbms.Set(2, "testing")
	dbms.Set(3, "testing")
	dbms.Set(4, "testing")
	dbms.Set(5, "testing")
	if dbms.Size() != 5 {
		t.Errorf("Wrong DBMS size after 5 inserts, %v instead of %v", dbms.Size(), 5)
	}
}

func TestDeleteOnlyCopy(t *testing.T) {
	dbms, _ := DBMS.Init(NonEmptyDBAbsPath)
	_ = dbms.Delete(4)
	val1, err1 := dbms.Get(4)
	if err1 == nil {
		t.Errorf("get after delete did not fail, val: %v, err: %v", val1, err1)
	}
}

func TestDeleteValWithMultipleVersions(t *testing.T) {
	dbms, _ := DBMS.Init(NonEmptyDBAbsPath)
	_ = dbms.Delete(2)
	val1, err1 := dbms.Get(2)
	if err1 == nil {
		t.Errorf("get after delete did not fail, val: %v, err: %v", val1, err1)
	}
}

func TestDeleteError(t *testing.T) {
	dbms, _ := DBMS.Init(NonEmptyDBAbsPath)
	err := dbms.Delete(1)
	if err == nil {
		t.Errorf("Did not receive an error for deleting non-existant ")
	}
}

func TestTruncate(t *testing.T) {
	dbms, _ := DBMS.Init(DbForSaveAbsPath)
	dbms.Set(1, "test")
	dbms.Set(2, "test2")
	dbms.Truncate()
	if dbms.Size() != 0 {
		t.Errorf("Truncate did not truncate hash table")
	}
	_ = dbms.Save()
	dbms, _ = DBMS.Init(DbForSaveAbsPath)
	if dbms.Size() != 0 {
		t.Errorf("Truncate did not clean the database")
	}
}

func TestSave(t *testing.T) {
	dbms, err := DBMS.Init(DbForSaveAbsPath)
	if err != nil {
		t.Fatal(err)
	}
	dbms.Set(1, "Test")
	dbms.Set(2, "Test2")
	err = dbms.Save()
	if err != nil {
		t.Fatal(err)
	}
	dbms, err = DBMS.Init(DbForSaveAbsPath)
	if err != nil {
		t.Fatal(err)
	}
	if dbms.Size() != 2 {
		t.Errorf("After saving and reading, size is %v instead of %v", dbms.Size(), 2)
	}
	val1, err1 := dbms.Get(1)
	if err1 != nil || val1 != "Test" {
		t.Errorf("Get key = %v failed, err: %v, val: %v", 1, err1, val1)
	}
	val2, err2 := dbms.Get(2)
	if err2 != nil || val2 != "Test2" {
		t.Errorf("Get key = %v failed, err: %v, val: %v", 2, err2, val2)
	}
	dbms.Truncate()
	_ = dbms.Save()
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
