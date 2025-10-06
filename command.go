package main

import (
	"github.com/google/shlex"
	"github.com/spf13/pflag"
)


var commands = make( map[string]*pflag.FlagSet )
var handlers = make( map[string]func(*pflag.FlagSet) )

type NormalDrop struct {
	
}

func (self *NormalDrop) String() string {
	return ""
}

func (self *NormalDrop) Set(s string) error {
	return nil
}

func (self *NormalDrop) Type() string {
	return ""
}



func overlay_tick() {
	
}



type Text struct {
	x, y int
	str string
	count int
}

func handleCommand_text(fs *pflag.FlagSet) {
	s := fs.Arg(0)
	if len(s) == 0 { return }
	printText(0, scrh+2, s)
}

func handleCommand_state(fs *pflag.FlagSet) {
	args := fs.Args()
	if len(args) == 0 {
		// maybe printout state info?
		return
	}
	switch args[0] {
		case "pause":
			rainStatus = !rainStatus
	}
}



func initCommands() {
	fset := pflag.NewFlagSet("text", pflag.ContinueOnError)
	fset.Int("xpos", 0, "x position of text")
	commands["text"] = fset
	handlers["text"] = handleCommand_text

	fset = pflag.NewFlagSet("state", pflag.ContinueOnError)
	commands["state"] = fset
	handlers["state"] = handleCommand_state
}


func processCommand(line string) {
	args, err := shlex.Split(line)
	if err != nil || len(args) == 0 { return }

	cmd, cmdExists := commands[args[0]]
	if !cmdExists { return }
	cmd.Parse(args[1:])

	hnd, hndExists := handlers[args[0]]
	if !hndExists { return }
	hnd(cmd)
}





type CellState struct {
	l0, l1, l2 bool
}

var state [][]CellState

func state_resize() {
	state = Reslice(state, scrw)
	for i := range state {
		state[i] = Reslice(state[i], scrh)
	}
}

func state_canDraw(x, y int, level int) bool {
	if x >= scrw { return false }
	if y >= scrh { return false }
	cell := &state[x][y]
	switch level {
		case 0:
			return cell.l1 | cell.l2
	}
	return false
}
