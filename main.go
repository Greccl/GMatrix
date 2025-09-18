package main

import (
	"fmt"
	"os"
	"time"
	"github.com/Greccl/tcell/v2"
)


type Color struct {
	r, g, b int32
}


var scr tcell.Screen
var w, h int
var headColor, tailColor Color
var rain Rain
var buffer [][]int8




func resize() {
	w, h = scr.Size()
	rain.resize(w, h)
	buffer = make([][]int8, w)
	for i:= range buffer {
		buffer[i] = make([]int8, h)
	}
}



func main() {
	headColor = Color{255, 153, 0}
	tailColor = Color{204, 51, 0}

	var e error
	scr, e = tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = scr.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defer scr.Fini()

	ch_ScreenEvents := make(chan tcell.Event)
	go func() {
		for {
			ev := scr.PollEvent()
			ch_ScreenEvents <- ev
		}
	}()
	ch_Tick := time.Tick(20*time.Millisecond)

	init := 100

	INIT: for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape {
							init = -2
							break INIT
						}
					case *tcell.EventResize:
						resize()
						init = 0
						break INIT
				}
			case <- ch_Tick:
				init--
				if init < 0 {
					break INIT
				}
		}
	}
	
	if init < 0 {
		return
	}

	LOOP: for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape { break LOOP }
					case *tcell.EventResize:
						resize()
				}
			case <- ch_Tick:
				rain.tick()
		}
	}
}
