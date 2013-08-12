package main

import (
	"gomux"
)

func main() {
	err := gomux.Init(5)
	if err != nil {
		panic(err)
	}
}
