package main

import (
	"fmt"
	"os"
	"time"
	// "bufio"
	"github.com/Greccl/tcell/v2"
)


type Color struct {
	r, g, b int32
}

var scr tcell.Screen
var scrw, scrh int
var oddOffset int


func resize() {
	scrw, scrh = scr.Size()
	scrh -= reservedHeight
	oddOffset = scrw % 2
	state_resize()
	rain_resize()
	scr.Clear()
}

func Reslice[T any](s []T, n int) []T {
	if n == len(s) { return s }
	if n <= cap(s) { return s[:n] }
	newSlice := make([]T, n)
	copy(newSlice, s)
	return newSlice
}

func blend(a, b Color, alfa int32) Color {
	var c Color
	beta := 1000 - alfa
	c.r = ((a.r * alfa) + (b.r * beta)) / 1000
	c.g = ((a.g * alfa) + (b.g * beta)) / 1000
	c.b = ((a.b * alfa) + (b.b * beta)) / 1000
	return c
}

func printText(x, y int, s string) {
	for i, r := range s {
		scr.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}
	scr.Show()
}

func drawCell(level, x, y int, r rune, s tcell.Style) {
	
}




func main() {
	// Read command line arguments
	defaults()
	initCommands()
	readCommandLine()
	
	// Open the pipe for reading external commands
	var ch_Commands chan string
	if cmdPath != "" {
		ch_Commands = make(chan string)
		go readCommandFile(cmdPath, ch_Commands)
	}

	// Init tcell screen
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

	// A timer to update animations
	ch_Tick := time.Tick(time.Duration(frameDuration)*time.Millisecond)

	// Setup process, wait for resize event or abort
	// if a timeout is reached (is 1 second enough?)
	initTimeout := 0

	INIT:
	for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape {
							initTimeout = -2
							break INIT
						}
					case *tcell.EventResize:
						resize()
						initTimeout = 0
						break INIT
				}
			case <- ch_Tick:
				initTimeout += frameDuration
				if initTimeout >= 1000 {
					break INIT
				}
		}
	}
	
	if initTimeout != 0 {
		// Timeout reached or aborted by Esc key
		return
	}

	LOOP:
	for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape { break LOOP }
						if ev.Key() == tcell.KeyRune {
							switch ev.Rune() {
								case 'p':
									rainStatus = !rainStatus
								case 's':
									if !rainStatus { rain_tick() }
								case 'q':
									break LOOP
							}
						}
					case *tcell.EventResize:
						resize()
				}
			case line := <- ch_Commands:
				processCommand(line)
				/*
				switch line {
					case "--end":
						break LOOP
					case "--pause":
						playState = !playState
				}
				*/
			case <- ch_Tick:
				if rainStatus {
					rain_tick()
				}
				overlay_tick()
		}
	}
}


