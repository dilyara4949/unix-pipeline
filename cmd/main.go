package main

import (
	"log"

	app "github.com/dilyara4949/unix-pipeline/internal/command"
)

func main() {
	cmds, out, err := app.ReadInput()
	if err != nil {
		log.Fatal(err.Error())
	}

	res, err := app.Execute(cmds, out)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, line := range res {
		log.Println(line)
	}
}
