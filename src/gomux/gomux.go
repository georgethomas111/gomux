package gomux

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"termbox-go"
	"unicode/utf8"
)

type Terminal struct {
	Width      int
	Height     int
	CursorX    int
	CursorY    int
	LinePrefix rune
	FgColor    termbox.Attribute
	BgColor    termbox.Attribute
	Stdout     io.ReadWriter
}

const (
	InputLength    = 100
	CommandLength  = 100
	DispLength     = 100
	DispRuneLength = 100
	OutFileName    = "/home/george/gomux.log"
)

var inputChan = make(chan termbox.Event, InputLength)
var processChan = make(chan string, CommandLength)
var dispChan = make(chan termbox.Event, DispLength)
var dispRuneChan = make(chan rune, DispRuneLength)
var drawSig = make(chan bool)

func getStdOut() (w io.ReadWriter) {

	w, err := os.Create(OutFileName)
	if err != nil {
		panic(err)
		return nil
	}
	return
}

func NewTerminal() (term *Terminal) {
	term = &Terminal{
		Width:      0,
		Height:     0,
		CursorX:    0,
		CursorY:    0,
		LinePrefix: '$',
		FgColor:    termbox.ColorDefault,
		BgColor:    termbox.ColorDefault,
		Stdout:     getStdOut(),
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
		var gomuxCom *exec.Cmd
		cmdArgs := strings.Split(command, " ")
		if len(cmdArgs) <= 1 {
			gomuxCom = exec.Command(command)
		} else {
			gomuxCom = exec.Command(cmdArgs[0],
				cmdArgs[1:len(cmdArgs)]...)
		}

		// Use custom stdout
		// Or have a worker which gets the contents of 
		// stdout and forwards it here
		gomuxCom.Stdout = t.Stdout
		gomuxCom.Stderr = t.Stdout
		gomuxCom.Stdin = os.Stdin
		gomuxCom.Run()
		drawSig <- true
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


// Use the t.Stdout read contents and draw on the screen.
func (t *Terminal) DrawFromFile() {

	for {
		<-drawSig
		fHandler, _ := os.Open(OutFileName)
		for {
			r := bufio.NewReader(fHandler)
			buf := make([]byte, 1024)
			read, err := r.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if read == 0 {
				break
			}
			fmt.Println(read)
			decodedCount := 0
			for {
				data, size := utf8.DecodeRune(buf)
				buf = buf[size:read]
				decodedCount += size
				dispRuneChan <- data
				if decodedCount == read {
					break
				}
			}
		}
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
	go t.DrawFromFile()
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
