package gomux

import (
	"fmt"
	"strconv"
	"termbox-go"
)

type Terminal struct {
	Width      int
	Height     int
	CursorX    int
	CursorY    int
	LinePrefix rune
	FgColor    termbox.Attribute
	BgColor    termbox.Attribute
}

const (
	InputLength    = 100
	CommandLength  = 100
	DispLength     = 100
	DispRuneLength = 100
)

var inputChan = make(chan termbox.Event, InputLength)
var processChan = make(chan string, CommandLength)
var dispChan = make(chan termbox.Event, DispLength)
var dispRuneChan = make(chan rune, DispRuneLength)

func NewTerminal() (term *Terminal) {
	term = &Terminal{
		Width:      0,
		Height:     0,
		CursorX:    0,
		CursorY:    0,
		LinePrefix: '$',
		FgColor:    termbox.ColorDefault,
		BgColor:    termbox.ColorDefault,
	}
	return
}

//
// initializes gomux with the number of panes to be shown in the output.
//
//
func (t *Terminal) Init() error {
	var err error
	fmt.Println("Starting gomux ... ")
	err = termbox.Init()
	if err != nil {
		return err
	}
	t.Width, t.Height = termbox.Size()
	defer termbox.Close()
	dispRuneChan <- t.LinePrefix
	err = t.Run()
	return err
}

func (t *Terminal) ProcessCommands() {
	for {
		command := <-processChan
		fmt.Println("Command", command)
	}
}

func (t *Terminal) DrawFromRune() {
	for {
		ch := <-dispRuneChan
		// Should a mutex lock be there ?
		termbox.SetCell(t.CursorX, t.CursorY,
			ch, t.FgColor, t.BgColor)
		t.CursorX++
		termbox.Flush()
	}
}

func (t *Terminal) DrawFromEvent() {
	for {
		event := <-dispChan
		switch event.Key {
		case termbox.KeyEnter:
			t.CursorY++
			t.CursorX = 0
			dispRuneChan <- t.LinePrefix
		default:
			t.CursorX++
			termbox.SetCell(t.CursorX, t.CursorY,
				event.Ch, t.FgColor, t.BgColor)
		}
		termbox.SetCursor(t.CursorX+1, t.CursorY)
		termbox.Flush()
	}
}

//
// Takes inputs and finally runs them when an enter key is hit.
//
//
func (t *Terminal) GetInput() {
	buffer := ""
	for {
		ch := ""
		event := <-inputChan
		switch event.Key {
		case termbox.KeyEnter:
			processChan <- buffer
			buffer = ""
		case termbox.KeySpace:
			ch = " "
		default:
			// TODO
			// 1. Decide on how to handle the errors.
			// 2. Decide on how to handle spaces.
			ch, _ = strconv.Unquote(strconv.QuoteRuneToASCII(event.Ch))
		}
		buffer += ch
	}
}

func (t *Terminal) Run() error {

	go t.GetInput()
	go t.ProcessCommands()
	go t.DrawFromEvent()
	go t.DrawFromRune()
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
				dispChan <- event
			}
		case termbox.EventError:
			return event.Err
		case termbox.EventResize:
			t.Width, t.Height = termbox.Size()
		}
	}
	return nil
}
