package gomux

import (
	"fmt"
	"termbox-go"
)

type RunDetails struct {

NoOfPanes int
Width int
Height int

}

var runDetails RunDetails
var inputChan = make(chan termbox.Key, 100)

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

//
// Takes inputs and finally runs them when an enter key is hit.
//
//

func GetInput() {

key := <-inputChan
switch key {
	case termbox.KeyEnter :
		fmt.Println("\n\n Enter key hit")
	}
}

func Run() error {

go GetInput()

loop:
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey :
			switch event.Key {
			case termbox.KeyEsc:
				break loop
			default:
				inputChan <- event.Key
			}
			case termbox.EventError :
				return event.Err
		}
	}
return nil
}
