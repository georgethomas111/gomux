package main

import (
	"gomux"
)

func main() {
	term := gomux.NewTerminal()
	err := term.Init()

	if err != nil {
		panic(err)
	}
}
