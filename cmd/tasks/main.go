package main

import (
	"log"
	"os"

	"github.com/EarvinKayonga/tasks/commands"
)

func main() {
	cmd := commands.Create()
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
