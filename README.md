gomux
=====

terminal multiplexer which executes same commands in different sessions

## Steps

1. Use termbox to get this done.
git clone https://github.com/nsf/termbox-go.git

2. Put the sources in such a way that GOPATH recognises it. 
Look at `go help path` for details.


## Input

1. The number of output sessions.
//TODO
2. Config file with number and location of the session to run in.
3. Have keystrokes to dynamically add/delete a session rather than from 
a config file.

## Design

1. Create the number of panes mentioned.
