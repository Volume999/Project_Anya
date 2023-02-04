package Client

import (
	"Project_Anya/Go_DB/DBMS"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func getInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func parseInput(tokens []string) (string, error) {
	command := tokens[0]
	switch command {
	case "help":
		if len(tokens) != 1 {
			return "", errors.New("help command does not take arguments")
		}
		return "Help Requested", nil
	case "get":
		if len(tokens) != 2 {
			return "", errors.New("invalid arguments for get, see help")
		}
		key := tokens[1]
		return fmt.Sprintf("Get Requested, Key = %v", key), nil
	case "set":
		if len(tokens) != 3 {
			return "", errors.New("invalid arguments for set, see help")
		}
		key := tokens[1]
		val := tokens[2]
		return fmt.Sprintf("Set Requested, Key = %v, Val = %v", key, val), nil
	default:
		return "", errors.New("invalid command")
	}
}

func Run() {
	reader := bufio.NewReader(os.Stdin)
	_, _ = DBMS.Init()
	//fmt.Println(db, err)
	for {
		input, _ := getInput(reader)
		tokens := strings.Fields(input)
		if len(tokens) != 0 {
			if output, err := parseInput(tokens); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(output)
			}
		}
	}

}
