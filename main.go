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

type Project struct{}

type Command struct {
	Name     string
	Argument string
}

var project Project

const (
	Cat  = "cat"
	Grep = "grep"
	Sort = "sort"
)

func (p *Project) Run() {
	cmds, out := readInput()
	res := execute(cmds, out)

	for _, line := range res {
		fmt.Println(line)
	}
}

func readInput() (<-chan Command, []string) {
	cmds := stdIn()

	var wg sync.WaitGroup
	wg.Add(1)

	var filePath string
	out := make(chan Command, len(cmds))

	go func(cmds []string, filePath *string) {
		for i, cmd := range cmds {
			sepCmd := strings.Fields(cmd)
			if i == 0 && len(sepCmd) == 1 {
				log.Fatal("filepath not given")
			}
			command := Command{}
			for k, name := range sepCmd {
				if command.Name == "" {
					command.Name = strings.ToLower(name)
				} else if name != "" {
					if *filePath == "" && k == len(sepCmd)-1 {
						*filePath = name
						wg.Done()
						continue
					} else {
						command.Argument = name
					}
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

func execute(cmds <-chan Command, input []string) []string {
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
	project.Run()
}
