package gomux

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

type Pane struct {

	Width      int
	Height     int
	CursorX    int
	CursorY    int
	InitX      int
	InitY      int
	LinePrefix rune
	FgColor    termbox.Attribute
	BgColor    termbox.Attribute
	Stdout     io.WriteCloser
	Stdin      io.ReadCloser
}

func NewPane(initX int, initY int, width int, height int) (pane *Pane) {

	reader, writer := getStdOut()
	pane = &Pane {
		Width:      width,
		Height:     height,
		InitX:      initX,
		InitY:      initY,
		CursorX:    initX,
		CursorY:    initY,
		LinePrefix: '$',
		FgColor:    termbox.ColorDefault,
		BgColor:    termbox.ColorDefault,
		Stdout:     writer,
		Stdin:      reader,
	}
	// Init the Pane
	pane.DrawFromRune(pane.LinePrefix)
	return
}


func (p *Pane) ProcessCommand(command string) {
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
		gomuxCom.Stdout = p.Stdout
		gomuxCom.Stderr = os.Stderr
		gomuxCom.Stdin = os.Stdin
		err := gomuxCom.Run()
		if(err != nil) {
			p.Stdout.Write([]byte(err.Error()))
		}
		p.DrawFromFile()
}

func (p *Pane) DrawFromRune(ch rune) {
		// Should a mutex lock be there ?
		if ch == '\n' {
			p.CursorX = p.InitX
			p.CursorY++
		} else {
			p.CursorX++
		}
		termbox.SetCell(p.CursorX, p.CursorY,
			ch, p.FgColor, p.BgColor)
		termbox.SetCursor(p.CursorX+1, p.CursorY)
		termbox.Flush()
}

func (p *Pane) DrawFromEvent(event termbox.Event) {
		switch event.Key {
		case termbox.KeyEnter:
			p.CursorY++
			p.CursorX = p.InitX
		default:
			p.CursorX++
			termbox.SetCell(p.CursorX, p.CursorY,
				event.Ch, p.FgColor, p.BgColor)
		}
		termbox.SetCursor(p.CursorX+1, p.CursorY)
		termbox.Flush()
}

// Use the t.Stdout read contents and draw on the screen.
func (p *Pane) DrawFromFile() {

		// Should be made equal to the maximum size of pipe
		// buffer wiki says unix one is 65536
		buf := make([]byte, 65536)
		read := -1
		read, err := p.Stdin.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		decodedCount := 0
		for decodedCount != (read - 1) {
			data, size := utf8.DecodeRune(buf)
			buf = buf[size:read]
			decodedCount += size
			p.DrawFromRune(data)
		}
		p.DrawFromRune('\n')
		p.DrawFromRune(p.LinePrefix)
}

