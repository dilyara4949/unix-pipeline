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
	Name     string
	Argument string
}

const (
	Cat  = "cat"
	Grep = "grep"
	Sort = "sort"
)

func init() {
	app = NewProject()
}

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
	cmds := stdIn()

	var wg sync.WaitGroup
	wg.Add(1)

	var filePath string
	out := make(chan Command, len(cmds))

	go func(cmds []string, filePath *string) {
		for _, cmd := range cmds {
			sepCmd := strings.Fields(cmd)
			if *filePath == "" && (len(sepCmd) == 1 || len(sepCmd) == 2 && sepCmd[0] == Grep) {
				log.Fatal("filepath not given")
			} else if len(sepCmd) == 0 {
				log.Fatal("input is not correct")
			}

			command := Command{}
			command.Name = strings.ToLower(sepCmd[0])
			if *filePath == "" {
				*filePath = sepCmd[len(sepCmd)-1]
				wg.Done()
			}
			if command.Name == Grep {
				if len(sepCmd) >= 2 {
					command.Argument = sepCmd[1]
				} else {
					log.Fatal("input is not correct")
				}
			}
			out <- command
		}
		close(out)
	}(cmds, &filePath)

	wg.Wait()
	fileText := readFile(filePath)
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

func stdIn() []string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(input), "|")
}

func readFile(path string) []string {
	dat, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(dat), "\n")
}

func main() {
	Run(app)
}
