package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

type project struct{}

type Project interface {
	ReadInput() (<-chan Command, []string)
	Execute(cmds <-chan Command, input []string) []string
	PrintResult(res []string)
}

var app Project

type Command struct {
	Name     Operation
	Argument string
}

type Operation string

const (
	Cat  Operation = "cat"
	Grep Operation = "grep"
	Sort Operation = "sort"
)

func NewProject() Project {
	return &project{}
}

func Run(p Project) {
	cmds, out := p.ReadInput()
	res := p.Execute(cmds, out)
	p.PrintResult(res)
}

func (p *project) PrintResult(res []string) {
	for _, line := range res {
		fmt.Println(line)
	}
}

func (p *project) ReadInput() (<-chan Command, []string) {
	cmds, err := stdIn()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var filePath string
	out := make(chan Command, len(cmds))

	go func(cmds []string) {
		for _, cmd := range cmds {
			sepCmd := strings.Fields(cmd)
			if len(sepCmd) == 0 {
				log.Fatal("input is not correct")
			}

			command := Command{}
			command.Name = Operation(strings.ToLower(sepCmd[0]))

			if len(sepCmd) < 2 && (filePath == "" || command.Name == Grep) {
				log.Fatal("input is not correct")
			}

			if command.Name == Grep {
				command.Argument = sepCmd[1]
			}

			if filePath == "" {
				filePath = sepCmd[len(sepCmd)-1]
				wg.Done()
			}
			out <- command
		}
		close(out)
	}(cmds)

	wg.Wait()
	fileText, err := readFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return out, fileText
}

func (p *project) Execute(cmds <-chan Command, input []string) []string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(cmds <-chan Command) {
		defer wg.Done()
		for cmd := range cmds {
			switch cmd.Name {
			case Cat:
				continue
			case Grep:
				grepFunc(&input, cmd.Argument)
			case Sort:
				sort.Strings(input)
			default:
				log.Fatal("unknown command")
			}
		}
	}(cmds)
	wg.Wait()
	return input
}

func grepFunc(input *[]string, arg string) {
	if arg == "" {
		log.Fatal("empty arg for grep")
	}
	for i := len(*input) - 1; i >= 0; i-- {
		if !strings.Contains((*input)[i], arg) {
			*input = append((*input)[:i], (*input)[i+1:]...)
		}
	}
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
	app = NewProject()
	Run(app)
}
