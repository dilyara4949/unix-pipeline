package command

import (
	"bufio"
	"fmt"
	"os"
	s "sort"
	"strings"
)

type Command struct {
	Name     Operation
	Argument string
}

type Operation string

const (
	cat            Operation = "cat"
	grep           Operation = "grep"
	sort           Operation = "sort"
	operationOrder int       = 0
	argumentOrder  int       = 1
)

func ReadInput() ([]Command, []string, error) {
	cmds, err := stdIn()
	if err != nil {
		return nil, nil, err
	}

	var filePath string

	out := make([]Command, len(cmds))

	for order, cmd := range cmds {
		sepCmd := strings.Fields(cmd)
		if len(sepCmd) == operationOrder {
			return nil, nil, fmt.Errorf("input is not correct")
		}

		c := Command{}
		c.Name = Operation(strings.ToLower(sepCmd[operationOrder]))

		if len(sepCmd) <= argumentOrder && (filePath == "" || c.Name == grep) {
			return nil, nil, fmt.Errorf("input is not correct")
		}

		if len(sepCmd) > argumentOrder {
			c.Argument = sepCmd[argumentOrder]
		}

		if filePath == "" {
			filePath = sepCmd[len(sepCmd)-1]
		}

		out[order] = c
	}

	fileText, err := readFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	return out, fileText, nil
}

func Execute(cmds []Command, input []string) ([]string, error) {
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
			input, err = grepFunc(input, cmd.Argument)
			if err != nil {
				return nil, err
			}
		case sort:
			s.Strings(input)
		default:
			return nil, fmt.Errorf("unknown command")
		}
	}

	return input, err
}

func grepFunc(input []string, arg string) ([]string, error) {
	if arg == "" {
		return nil, fmt.Errorf("empty arg for grep")
	}

	for i := len(input) - 1; i >= 0; i-- {
		if !strings.Contains(input[i], arg) {
			input = append(input[:i], input[i+1:]...)
		}
	}

	return input, nil
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
