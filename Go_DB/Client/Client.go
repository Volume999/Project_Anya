package Client

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := getInput(reader)
		fmt.Println(input)
	}
}
