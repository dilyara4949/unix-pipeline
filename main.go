package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type command struct {
	name     operation
	argument string
}

type operation string

const (
	cat            operation = "cat"
	grep           operation = "grep"
	ssort          operation = "sort"
	operationOrder int       = 0
	argumentOrder  int       = 1
)

func readInput() ([]command, []string, error) {
	cmds, err := stdIn()
	if err != nil {
		return nil, nil, err
	}

	var filePath string
	out := make([]command, len(cmds))

	for i, cmd := range cmds {
		sepCmd := strings.Fields(cmd)
		if len(sepCmd) == operationOrder {
			return nil, nil, fmt.Errorf("input is not correct")
		}

		c := command{}
		c.name = operation(strings.ToLower(sepCmd[operationOrder]))

		if len(sepCmd) <= argumentOrder && (filePath == "" || c.name == grep) {
			return nil, nil, fmt.Errorf("input is not correct")
		}

		if len(sepCmd) > argumentOrder {
			c.argument = sepCmd[argumentOrder]
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

func execute(cmds []command, input []string) ([]string, error) {
	var err error
	for _, cmd := range cmds {
		switch cmd.name {
		case cat:
			if cmd.argument != "" {
				input, err = readFile(cmd.argument)
				if err != nil {
					return nil, err
				}
			}
		case grep:
			err = grepFunc(&input, cmd.argument)
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

func main() {
	cmds, out, err := readInput()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	res, err := execute(cmds, out)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	for _, line := range res {
		fmt.Println(line)
	}
}
