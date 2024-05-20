package internal

import (
	"bufio"
	"fmt"
	app "github.com/dilyara4949/unix-pipeline"
	"os"
	"sort"
	"strings"
)

const (
	cat            app.Operation = "cat"
	grep           app.Operation = "grep"
	ssort          app.Operation = "sort"
	operationOrder int           = 0
	argumentOrder  int           = 1
)

func ReadInput() ([]app.Command, []string, error) {
	cmds, err := stdIn()
	if err != nil {
		return nil, nil, err
	}

	var filePath string
	out := make([]app.Command, len(cmds))

	for i, cmd := range cmds {
		sepCmd := strings.Fields(cmd)
		if len(sepCmd) == operationOrder {
			return nil, nil, fmt.Errorf("input is not correct")
		}

		c := app.Command{}
		c.Name = app.Operation(strings.ToLower(sepCmd[operationOrder]))

		if len(sepCmd) <= argumentOrder && (filePath == "" || c.Name == grep) {
			return nil, nil, fmt.Errorf("input is not correct")
		}

		if len(sepCmd) > argumentOrder {
			c.Argument = sepCmd[argumentOrder]
		}

		if filePath == "" {
			filePath = sepCmd[len(sepCmd)-1]
		}
		out[i] = c
	}
	fileText, err := readFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	return out, fileText, nil
}

func Execute(cmds []app.Command, input []string) ([]string, error) {
	var err error
	for _, cmd := range cmds {
		switch cmd.Name {
		case cat:
			if cmd.Argument != "" {
				input, err = readFile(cmd.Argument)
				if err != nil {
					return nil, err
				}
			}
		case grep:
			err = grepFunc(&input, cmd.Argument)
			if err != nil {
				return nil, err
			}
		case ssort:
			sort.Strings(input)
		default:
			return nil, fmt.Errorf("unknown command")
		}
	}
	return input, err
}

func grepFunc(input *[]string, arg string) error {
	if arg == "" {
		return fmt.Errorf("empty arg for grep")
	}
	for i := len(*input) - 1; i >= 0; i-- {
		if !strings.Contains((*input)[i], arg) {
			*input = append((*input)[:i], (*input)[i+1:]...)
		}
	}
	return nil
}

func stdIn() ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(input), "|"), nil
}

func readFile(path string) ([]string, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(dat), "\n"), nil
}