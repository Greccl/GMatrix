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
}



func overlay_touch(x, y int) {
	for _, t := range boxes {
		if y == t.y {
			// t := boxes[i]
			l := len(t.runes)
			p := x - t.x
			if p >= 0 && p < l {
				t.blends[p] = 1000
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

	id, _ := fs.GetInt("id")
	if id > 0 {
		t = findTextById(id)
	}
	name, _ := fs.GetString("name")
	if t == nil && name != "" {
		t = findTextByName(name)
	}

	if b, _ := fs.GetBool("kill"); b {
		if t != nil {
			t.autoremove()
		}
		return
	}

	if t == nil {
		t = NewText()
		t.id = id
		t.name = name
		boxes = append(boxes, t)
	}

	s := fs.Arg(0)
	if len(s) > 0 {
		t.setText(s)
	}

	s, _ = fs.GetString("foreground")
	if len(s) > 0 {
		c, err := parseColor(s)
		if err == nil { t.fg = c }
		// {panic("FG NOT SET")}
		// panic("FG SET")
	}

	s, _ = fs.GetString("background")
	if len(s) > 0 {
		c, err := parseColor(s)
		if err == nil { t.bg = c }
		// panic("SET BG")
	}

	if b, _ := fs.GetBool("centerx"); b {
		t.centerx = true
	}
	if b, _ := fs.GetBool("centery"); b {
		t.centery = true
	}

	var moved bool
	x, err := fs.GetInt("x")
	if err == nil { moved = true }
	y, err := fs.GetInt("y")
	if err == nil { moved = true }
	if moved { t.move(x, y) }

	s, _ = fs.GetString("animation")
	if len(s) > 0 {
		t.setAnimation(s)
	} else {
		if t.animName != "" {
			t.animFunc()
		} else {
			// b := t.animData_bool
			t.animData_bool = false
			t.anim_0()
			// t.animData_bool = b
		}
	}

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
	fset.IntP    ("x"         , "x", 0     , "x position of text box")
	fset.IntP    ("y"         , "y", 0     , "y position of text box")
	fset.BoolP   ("centerx"   , "X", false , "evaluate horizontal position from center of screen")
	fset.BoolP   ("centery"   , "Y", false , "evaluate vertical position from center of screen")
	fset.IntP    ("id"        , "i", 0     , "identifier (integer value) for the text box")
	fset.StringP ("name"      , "n", ""    , "identifier (string) for the text box")
	fset.BoolP   ("modify"    , "m", false , "try modify an existent box")
	fset.BoolP   ("kill"      , "k", false , "try remove given box")
	fset.StringP ("animation" , "a", "", "name of animation")
	fset.StringP ("foreground", "f", ""    , "set foreground colour")
	fset.StringP ("background", "b", ""    , "set background colour")

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
