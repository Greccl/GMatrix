package main

import (
	// "math/rand/v2"
	"github.com/google/shlex"
	"github.com/spf13/pflag"
	// "github.com/Greccl/tcell/v2"
)





var commands = make( map[string]*pflag.FlagSet )
var handlers = make( map[string]func(*pflag.FlagSet) )
var overlayDivider int
var boxes []*Text


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





func overlay_touch(x, y int) {
	for _, t := range boxes {
		if y == t.y {
			// t := boxes[i]
			l := len(t.runes)
			p := x - t.x
			if p >= 0 && p < l {
				// t.blends[p] = 1000
				break
			}
		}
	}
}

func overlay_tick() {
/*
	if overlayDivider < 2 {
		overlayDivider++
		return
	}
	overlayDivider = 0
*/	
	for i := range boxes {
		if boxes[i] == nil { continue }
		t := boxes[i]
		// t.animFunc()
		if t.anim == nil { continue }
		t.anim.tick(t)
	}
}




func findTextById(id int) *Text {
	for i := range boxes {
		if boxes[i].id == id {
			return boxes[i]
		}
	}
	return nil
}

func findTextByName(name string) *Text {
	for i := range boxes {
		if boxes[i].name == name {
			return boxes[i]
		}
	}
	return nil
}

func handleCommand_text(fs *pflag.FlagSet) {
	var t *Text
	var draw, movex, movey bool

	// Find target text
	id, _ := fs.GetInt("id")
	if id > 0 {
		t = findTextById(id)
	}
	name, _ := fs.GetString("name")
	if t == nil && name != "" {
		t = findTextByName(name)
	}

	// Issue a kill command
	if b, _ := fs.GetBool("kill"); b {
		if t != nil {
			t.autoremove()
		}
		return
	}

	// Create if doesnt exists
	if t == nil {
		t = NewText()
		t.id = id
		t.name = name
		boxes = append(boxes, t)
	}

	// Text
	s := fs.Arg(0)
	if len(s) > 0 {
		t.setText(s)
		draw = true
	}

	// Colours
	s, _ = fs.GetString("foreground")
	if len(s) > 0 {
		c, err := parseColor(s)
		if err == nil {
			t.fg = c
			draw = true
		}
	}

	s, _ = fs.GetString("background")
	if len(s) > 0 {
		c, err := parseColor(s)
		if err == nil {
			t.bg = c
			draw = true
		}
	}

	// Align
	if changed, value := getBool(fs, "halign"); changed {
		t.halign = value
		movex = true
	}
	if changed, value := getBool(fs, "valign"); changed {
		t.valign = value
		movey = true
	}

	// Position
	if changed, value := getIntSlice(fs, "position"); changed {
		if len(value) == 2 {
			t.setx = value[0]
			t.sety = value[1]
			movex = true
			movey = true
		}
	} else {
		if changed, value := getInt(fs, "x"); changed {
			t.setx = value
			movex = true
		}
		if changed, value := getInt(fs, "y"); changed {
			t.sety = value
			movey = true
		}
	}

	// Animation type
	if changed, value := getString(fs, "animation"); changed {
		t.setAnimation(value)
		draw = true
	}

	// recalculate position
	if movex {
		t.movex()
		draw = true
	}

	if movey {
		t.movey()
		draw = true
	}

	// redraw needed
	if draw {
		t.anim.draw(t)
	}

	// Reset flagset state to process next command
	fs.VisitAll(resetFlag)
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
	commands["text"] = fset
	handlers["text"] = handleCommand_text
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

	fset = pflag.NewFlagSet("state", pflag.ContinueOnError)
	commands["state"] = fset
	handlers["state"] = handleCommand_state

/*
	fset = pflag.NewFlagSet("delay", pflag.ContinueOnError)
	commands["automate"] = fset
	handlers["state"] = handleCommand_state
*/
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
