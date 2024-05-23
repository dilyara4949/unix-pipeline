package main

import (
	"log"

	c "github.com/dilyara4949/unix-pipeline/internal/command"
)

func main() {
	cmds, out, err := c.ReadInput()
	if err != nil {
		log.Fatal(err.Error())
	}

	res, err := c.Execute(cmds, out)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, line := range res {
		log.Println(line)
	}
}
