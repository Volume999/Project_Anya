package Evaluation

import (
	"Project_Anya/GoDB/DBMS"
	"fmt"
	"github.com/dchest/uniuri"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	LoadUpperBound          = 100000
	LoadLowerBound          = 10000
	LoadIncrement           = 10000
	TestNumberOfRepetitions = 10000
	StatsFilePath           = "./GoDB/Evaluation/Report/GoDBEvalStats.csv"
	ReadDBPath              = "./GoDB/Evaluation/Mock_DBs/GoDBEvalReadDB"
	WriteDBPath             = "./GoDB/Evaluation/Mock_DBs/GoDBEvalWriteDB"
	LargeStringLength       = 100
)

var statsFile *os.File
var readDB DBMS.Dbms
var writeDB DBMS.Dbms

func genString(size int) string {
	return uniuri.NewLen(size)
}

func setupReadDb() {
	for i := 0; i < LoadUpperBound; i++ {
		str := genString(LargeStringLength)
		readDB.Set(i, str)
	}
	err := readDB.Save()
	if err != nil {
		panic("could not save dbms for setup")
	}
}

//	func record_stats(row string) {
//		statsAbsPath, err := filepath.Abs(StatsFilePath)
//		os.OpenFile(statsAbsPath)
//		os.WriteFile(statsAbsPath, row)
//	}
func test_get(inputSize int) {
	for i := 0; i < TestNumberOfRepetitions; i++ {
		id := rand.Intn(inputSize)
		_, _ = readDB.Get(id)
	}
}

func test_set(inputSize int) {
	for i := 0; i < inputSize; i++ {
		val := genString(LargeStringLength)
		writeDB.Set(i, val)
	}
}

func document(begin time.Time) {

}

func test() {
	defer document(time.Now())
}

func init() {
	statsAbsPath, err := filepath.Abs(StatsFilePath)
	if err = os.Truncate(statsAbsPath, 0); err != nil {
		panic(err)
	}
	//header := []string{"TestName", "InputSize", "Runtime(seconds)", "Latency", "Throughput"}
	header := "TestName,InputSize,Runtime(seconds),Latency,Throughput\n"
	err = os.WriteFile(statsAbsPath, []byte(header), 0755)
	statsFile, _ = os.OpenFile(statsAbsPath, os.O_APPEND, 0755)
	evalReadDBPath, _ := filepath.Abs(ReadDBPath)
	readDB, err = DBMS.Init(evalReadDBPath)
	if err != nil {
		panic(err)
	}
	if readDB.Size() == 0 {
		setupReadDb()
	}
	evalWriteDBPath, _ := filepath.Abs(WriteDBPath)
	writeDB, _ = DBMS.Init(evalWriteDBPath)
}

func Eval() {
	fmt.Println("Running")
	defer func() {
		err := statsFile.Close()
		if err != nil {
			panic("could not close stats file")
		}
	}()
}
