package Client

import (
	"Project_Anya/GoDB/DBMS"
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Client struct {
	dbms   *DBMS.Dbms
	reader *bufio.Reader
}

func Init(dbms *DBMS.Dbms, reader *bufio.Reader) Client {
	return Client{dbms: dbms, reader: reader}
}

func (client *Client) getInput() (string, error) {
	input, err := client.reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func (client *Client) parseInput(input string) (string, error) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return "", nil // Empty string as command - ignore
	}
	out, err := client.parseCommand(tokens)
	return out, err
}

func (client *Client) parseCommand(tokens []string) (string, error) {
	command := tokens[0]
	switch command {
	case "help":
		if len(tokens) != 1 {
			return "", errors.New("help command does not take arguments")
		}
		output := `API for AnyaDB:
get {key: int} -> returns key if found, else returns an error
set {key: int} {val: int} -> inserts value in db
del {key: int} -> deletes the value from the db
save -> save db state
exit -> exit client
 `
		return output, nil
	case "get":
		if len(tokens) != 2 {
			return "", errors.New("invalid arguments for get, see help")
		}
		if key, err := strconv.Atoi(tokens[1]); err == nil {
			return client.dbms.Get(key)
		} else {
			return "", errors.New("key must be an integer")
		}
	case "set":
		if len(tokens) != 3 {
			return "", errors.New("invalid arguments for set, see help")
		}
		key, err := strconv.Atoi(tokens[1])
		if err != nil {
			return "", errors.New("key must be an integer")
		}
		val := tokens[2]
		client.dbms.Set(key, val)
		return fmt.Sprintf("Set Requested, Key = %v, Val = %v", key, val), nil
	case "save":
		if len(tokens) != 1 {
			return "", errors.New("save does not take parameters")
		}
		_ = client.dbms.Save()
		return "", nil
	case "exit":
		if len(tokens) != 1 {
			return "", errors.New("exit does not take parameters")
		}
		return "", errors.New("exit")
	case "del":
		if len(tokens) != 2 {
			return "", errors.New("invalid parameters for delete, see help")
		}
		if key, err := strconv.Atoi(tokens[1]); err != nil {
			return "", errors.New("key must be an integer")
		} else {
			if err := client.dbms.Delete(key); err != nil {
				return "", errors.New("key not found")
			} else {
				return "ok", nil
			}
		}

	default:
		return "", errors.New("invalid command")
	}
}

func (client *Client) Run() {
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Println(db, err)
	for {
		input, _ := client.getInput()
		output, err := client.parseInput(input)
		if err == nil {
			fmt.Println(output)
		} else {
			if err.Error() == "exit" {
				if err := client.dbms.Save(); err != nil {
					fmt.Println("Failed to save database")
				}
				break
			}
			fmt.Println(err)
		}
	}
}
