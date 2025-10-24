package main

import (
	"github.com/google/shlex"
	"github.com/spf13/pflag"
)





type Command struct {
	fs *pflag.FlagSet
	fn func(*pflag.FlagSet)
}

type SplitCmd struct {
	cmd *Command
	args []string
}

var commands = make( map[string]*Command )





//
// Helpers for FlagSet manipulation
//

func resetFlag(f *pflag.Flag) {
	f.Value.Set(f.DefValue)
	f.Changed = false
}

func getBool(fs *pflag.FlagSet, name string) (bool,bool) {
	f := fs.Lookup(name)
	if !f.Changed { return false, false }
	value, _ := fs.GetBool(name)
	return true, value
}

func getInt(fs *pflag.FlagSet, name string) (bool,int) {
	f := fs.Lookup(name)
	if !f.Changed { return false, 0 }
	value, _ := fs.GetInt(name)
	return true, value
}

func getIntSlice(fs *pflag.FlagSet, name string) (bool,[]int) {
	f := fs.Lookup(name)
	if !f.Changed { return false, []int{} }
	value, _ := fs.GetIntSlice(name)
	return true, value
}

func getString(fs *pflag.FlagSet, name string) (bool,string) {
	f := fs.Lookup(name)
	if !f.Changed { return false, "" }
	value, _ := fs.GetString(name)
	return true, value
}
















func handleCommand_state(fs *pflag.FlagSet) string {
	args := fs.Args()
	if len(args) == 0 {
		// maybe printout state info?
		return ""
	}
	switch args[0] {
		case "pause":
			rainStatus = !rainStatus
	}
	return ""
}





func initCommands() {
	fset := pflag.NewFlagSet("text", pflag.ContinueOnError)
	fset.IntSliceP("position" , "p", []int{}, "a compact way to set x and y position")
	fset.IntP    ("x"         , "x", 0      , "x position of text box")
	fset.IntP    ("y"         , "y", 0      , "y position of text box")
	fset.BoolP   ("halign"    , "h", false  , "evaluate horizontal position from center of screen")
	fset.BoolP   ("valign"    , "v", false  , "evaluate vertical position from center of screen")
	fset.IntP    ("id"        , "i", 0      , "identifier (integer value) for the text box")
	fset.StringP ("name"      , "n", ""     , "identifier (string) for the text box")
	fset.BoolP   ("kill"      , "k", false  , "try remove given box")
	fset.StringP ("animation" , "a", ""     , "name of animation")
	fset.StringP ("foreground", "f", ""     , "set foreground colour")
	fset.StringP ("background", "b", ""     , "set background colour")
	commands["text"] = &Command{fset,handleCommand_text}
	// cmd := new(Command)
	// commands["text"] = cmd
	// cmd.fs = fset
	// cmd.fn = handleCommand_text

	// fset = pflag.NewFlagSet("state", pflag.ContinueOnError)
	// commands["state"] = fset
	// handlers["state"] = handleCommand_state

/*
	fset = pflag.NewFlagSet("delay", pflag.ContinueOnError)
	commands["automate"] = fset
	handlers["state"] = handleCommand_state
*/
}






func parseCommand(line string) {
	args, err := shlex.Split(line)
	if err != nil || len(args) == 0 { return }
	cmd, exists := commands[args[0]]
	if !exists { return }
	if len(args) > 0 {
		ch_Commands <- SplitCmd{cmd,args[1:]}
	} else {
		ch_Commands <- SplitCmd{cmd,args[1:]}
	}
}

