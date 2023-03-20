package Evaluation

import (
	"Project_Anya/GoDB/DBMS"
	"fmt"
	"github.com/dchest/uniuri"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	LoadUpperBound          = 200000
	LoadLowerBound          = 10000
	LoadIncrement           = 10000
	TestNumberOfRepetitions = 10000
	StatsFilePath           = "./GoDB/Evaluation/Report/GoDBEvalStats.csv"
	ReadDBPath              = "./GoDB/Evaluation/Mock_DBs/GoDBEvalReadDB"
	WriteDBPath             = "./GoDB/Evaluation/Mock_DBs/GoDBEvalWriteDB"
	InitDBPath              = "./GoDB/Evaluation/Mock_DBs/GoDBEvalInitDB"
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
func testGet(inputSize int) {
	for i := 0; i < TestNumberOfRepetitions; i++ {
		id := rand.Intn(inputSize)
		_, _ = readDB.Get(id)
	}
}

func testSet(inputSize int) {
	for i := 0; i < inputSize; i++ {
		val := genString(LargeStringLength)
		writeDB.Set(i, val)
	}
}

func testInitPrepare(inputSize int) {
	evalReadDBPath, _ := filepath.Abs(ReadDBPath)
	data, _ := os.ReadFile(evalReadDBPath)
	evalInitDBPath, _ := filepath.Abs(InitDBPath)
	err := os.WriteFile(evalInitDBPath, data, 0755)
	if err != nil {
		panic(fmt.Sprintf("Test Init fail, error: %v", err))
	}
	initDB, err := DBMS.Init(evalInitDBPath)
	if err != nil {
		panic(fmt.Sprintf("Test Init fail, error: %v", err))
	}
	for i := inputSize + 1; ; i++ {
		if err := initDB.Delete(i); err != nil {
			break
		}
	}
	err = initDB.Save()
	if err != nil {
		panic(fmt.Sprintf("Test Init fail, error: %v", err))
	}
}

func testInit(_ int) {
	evalInitDBPath, _ := filepath.Abs(InitDBPath)
	_, err := DBMS.Init(evalInitDBPath)
	if err != nil {
		panic(fmt.Sprintf("Test Init fail, error: %v", err))
	}
}

func getFunctionName(temp interface{}) string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name(), ".")
	return strs[len(strs)-1]
}

func document(testName string, inputSize int, begin time.Time) {
	elapsedTimeInSeconds := time.Since(begin).Seconds()
	throughput := float64(inputSize) / elapsedTimeInSeconds
	latency := elapsedTimeInSeconds / float64(inputSize)
	if _, err := statsFile.WriteString(fmt.Sprintf("%v,%v,%.3f,%.5f,%.3f\n", testName, inputSize, elapsedTimeInSeconds, latency, throughput)); err != nil {
		panic(fmt.Sprintf("could not write results of test, err: %v", err))
	}
}

func test(testFunc func(int), inputSize int) {
	functionName := getFunctionName(testFunc)
	if functionName == "testGet" {
		writeDB.Truncate()
	}
	if functionName == "testInit" {
		testInitPrepare(inputSize)
	}
	beginTime := time.Now()
	defer document(functionName, inputSize, beginTime)
	testFunc(inputSize)
}

func init() {
	statsAbsPath, err := filepath.Abs(StatsFilePath)
	if err = os.Truncate(statsAbsPath, 0); err != nil {
		panic(err)
	}
	//header := []string{"TestName", "InputSize", "Runtime(seconds)", "Latency", "Throughput"}
	header := "TestName,InputSize,Runtime(seconds),Latency,Throughput\n"
	err = os.WriteFile(statsAbsPath, []byte(header), 0755)
	statsFile, _ = os.OpenFile(statsAbsPath, os.O_WRONLY|os.O_APPEND, 0755)
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

func RunEval() {
	defer func() {
		err := statsFile.Close()
		if err != nil {
			panic("could not close stats file")
		}
		initDBAbsPath, _ := filepath.Abs(InitDBPath)
		err = os.Remove(initDBAbsPath)
		if err != nil {
			panic("Clean up error: could not delete init db")
		}
	}()
	fmt.Println("Running")
	for inputSize := LoadLowerBound; inputSize <= LoadUpperBound; inputSize += LoadIncrement {
		test(testGet, inputSize)
	}
	for inputSize := LoadLowerBound; inputSize <= LoadUpperBound; inputSize += LoadIncrement {
		test(testSet, inputSize)
	}
	for inputSize := LoadLowerBound; inputSize <= LoadUpperBound; inputSize += LoadIncrement {
		test(testInit, inputSize)
	}
}
