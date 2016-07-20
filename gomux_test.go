package gomux

import (
	"testing"
)

func TestMain(m *testing.M) {
	term := NewTerminal()
	err := term.Init()

	if err != nil {
		panic(err)
	}
}
