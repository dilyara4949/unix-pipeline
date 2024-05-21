package main

import (
	"log"
	"os"

	app "github.com/dilyara4949/unix-pipeline/internal"
)

func main() {
	cmds, out, err := app.ReadInput()
	if err != nil {
		log.Println(err.Error())
		os.Exit(0)
	}

	res, err := app.Execute(cmds, out)
	if err != nil {
		log.Println(err.Error())
		os.Exit(0)
	}

	for _, line := range res {
		log.Println(line)
	}
}
