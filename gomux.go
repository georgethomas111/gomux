package gomux

import (
	"fmt"
	"io"
	"os"
	"github.com/nsf/termbox-go"
	"strconv"
)

type Terminal struct {
	Width    int
	Height   int
	Panes      []*Pane
}

const (
	InputLength    = 100
	CommandLength  = 100
	DispLength     = 100
	DispRuneLength = 100
)

	var InputChan = make(chan termbox.Event, InputLength)
	var ProcessChan = make(chan string, CommandLength)
	var DispChan = make(chan termbox.Event, DispLength)
	var DispRuneChan = make(chan rune, DispRuneLength)
	var DrawSig = make(chan bool)

func getStdOut() (rFile io.ReadCloser, wFile io.WriteCloser) {

	rFile, wFile, err := os.Pipe()
	if err != nil {
		panic(err)
		return nil, nil
	}
	return
}

func NewTerminal() (term *Terminal) {

	term = &Terminal{
		Width:      0,
		Height:     0,
	}
	return
}

//
// initializes gomux with the number of panes to be shown in the output.
// Add one pane to it by default.
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
	pane := NewPane(0, 0, t.Width/2, t.Height/2)
	t.Panes = append(t.Panes, pane)
	pane = NewPane(t.Width/2, 0, t.Width/2, t.Height)
	t.Panes = append(t.Panes, pane)
	pane = NewPane(0, t.Height/2, t.Width/2, t.Height/2)
	t.Panes = append(t.Panes, pane)
	err = t.Run()
	return err
}

func (t* Terminal) GetInput() {

	buffer := ""
	for {
		ch := ""
		event := <-InputChan
		switch event.Key {
		case termbox.KeyEnter:
			ProcessChan <- buffer
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

func (t* Terminal) ProcessCommands() {
	for {
		command := <-ProcessChan
		for _, pane := range t.Panes {
			pane.ProcessCommand(command)
		}

	}
}

func (t* Terminal) DrawFromEvent() {
	for {
		event := <-DispChan
		for _, pane := range t.Panes {
			pane.DrawFromEvent(event)
		}
	}
}

func (t* Terminal) DrawFromRune() {
	for {
		ch := <-DispRuneChan
		for _, pane := range t.Panes {
			pane.DrawFromRune(ch)
		}
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
				InputChan <- event
				DispChan <- event
			}
		case termbox.EventError:
			return event.Err
		case termbox.EventResize:
			t.Width, t.Height = termbox.Size()
		}
	}
	return nil
}
