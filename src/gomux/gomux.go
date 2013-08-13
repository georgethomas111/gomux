package gomux

import (
	"fmt"
	"termbox-go"
	"strconv"
)

type Terminal struct {
	Width      int
	Height     int
        CursorX    int
        CursorY    int
	LinePrefix string
	buffer     string
}

var inputChan = make(chan termbox.Event, 100)
var processChan = make(chan string, 100)
var dispChan = make(chan string, 100)

func NewTerminal() (term *Terminal) {
	term = &Terminal {
		Width : 0,
		Height : 0,
		CursorX : 0,
		CursorY : 0,
		LinePrefix : "$",
	}
	return
}

//
// initializes gomux with the number of panes to be shown in the output.
//
//
func(t* Terminal) Init() error {
	var err error
	fmt.Println("Starting gomux ... ")
	err = termbox.Init()
	if err != nil {
		return err
	}
	t.Width, t.Height = termbox.Size()
	defer termbox.Close()
	err = t.Run()
	return err
}


func(t* Terminal) ProcessCommands() {
	for {
		command := <-processChan
		fmt.Println(command)
	}
}


func(t* Terminal) Draw() {
	for {
		char := <-dispChan
		t.CursorX += len(char)
		// TODO Create a cross platform method for endline
		if char == "\n" {
			t.CursorY += 1
			t.buffer  = t.LinePrefix
			t.CursorX = len(t.buffer)
		}
		fmt.Println(t.CursorX,":",t.CursorY)
	}
}

//
// Takes inputs and finally runs them when an enter key is hit.
//
//
func(t* Terminal) GetInput() {
	buffer := ""
	for {
		ch := ""
		event := <-inputChan
		switch event.Key {
		case termbox.KeyEnter:
			processChan <- buffer
			buffer  = ""
			dispChan <- "\n"
			continue
		case termbox.KeySpace:
			ch = " "
		default:
			// TODO 
			// 1. Decide on how to handle the errors.
			// 2. Decide on how to handle spaces.
			ch, _  = strconv.Unquote(strconv.QuoteRuneToASCII(event.Ch))
		}
		buffer += ch
		dispChan <- ch
	}
}

func(t* Terminal) Run() error {

	go t.GetInput()
	go t.ProcessCommands()
        go t.Draw()
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
		case termbox.EventResize:
			t.Width, t.Height = termbox.Size()
		}
	}
	return nil
}
