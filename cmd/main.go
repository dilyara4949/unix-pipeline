package main

import (
	"fmt"
	app "github.com/dilyara4949/unix-pipeline/internal"
	"os"
)

func main() {
	cmds, out, err := app.ReadInput()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	res, err := app.Execute(cmds, out)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	for _, line := range res {
		fmt.Println(line)
	}
}
