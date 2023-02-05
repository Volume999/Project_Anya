package Testing

import (
	"Project_Anya/Go_DB/DBMS"
	"fmt"
	"path/filepath"
)

const (
	DBABSPATH = "Go_DB/Database/db"
)

func initDBMS() DBMS.Dbms {
	dbPath, _ := filepath.Abs(DBABSPATH)
	dbms, _ := DBMS.Init(dbPath)
	return dbms
}

func testGet(dbms DBMS.Dbms, key int) {
	val, err := dbms.Get(key)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(val)
	}
}

func TestDBMS() {
	dbms := initDBMS()
	dbms.Set(4, "Alilol")
	dbms.Set(2, "AliPizza")
	dbms.Set(3, "AliAli")
	testGet(dbms, 4)
	testGet(dbms, 3)
	testGet(dbms, 1)
	_ = dbms.Delete(3)
	_ = dbms.Delete(4)
	testGet(dbms, 2)
	testGet(dbms, 3)
	testGet(dbms, 4)
	_ = dbms.Save()
}
