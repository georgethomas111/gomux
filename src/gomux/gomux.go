package gomux

import (
	"fmt"
	"termbox-go"
	"strconv"
)

type RunDetails struct {
	NoOfPanes int
	Width     int
	Height    int
}

var runDetails RunDetails
var inputChan = make(chan termbox.Event, 100)
var processChan = make(chan string, 100)

//
// initializes gomux with the number of panes to be shown in the output.
//
//

func Init(noPanes int) error {
	var err error
	runDetails.NoOfPanes = noPanes
	fmt.Println("Starting gomux ... ")
	err = termbox.Init()
	if err != nil {
		return err
	}
	runDetails.Width, runDetails.Height = termbox.Size()
	defer termbox.Close()
	err = Run()
	return err
}


func ProcessCommands() {
	for {
		command := <-processChan
		fmt.Println("$" + command)
	}
}

//
// Takes inputs and finally runs them when an enter key is hit.
//
//

func GetInput() {
	buffer := ""
	for {
		event := <-inputChan
		switch event.Key {
		case termbox.KeyEnter:
			processChan <- buffer
			buffer = ""
		default:
			// TODO 
			// 1. Decide on how to handle the errors.
			// 2. Decide on how to handle spaces.
			ch, _  := strconv.Unquote(strconv.QuoteRuneToASCII(event.Ch))
			buffer += ch
		}
	}
}

func Run() error {

	go GetInput()
	go ProcessCommands()
loop:
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyEsc:
				break loop
			default:
				inputChan <- event
			}
		case termbox.EventError:
			return event.Err
		}
	}
	return nil
}
