package main

import (
	"fmt"
	"termbox-go"
)

func main() {
	fmt.Println("Starting gomux ... ")
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

loop:
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeyEsc:
				break loop
			default:
				fmt.Println("recieved key stroke", string(event.Key))
			}
		}
	}
}
