package main

import (
	"os"
	"bufio"
	"github.com/spf13/pflag"
)
// import flags "github.com/jessevdk/go-flags"

var lucentHead, lucentBody Color
var backHead, backNeck, backTail Color

var normalHead, normalNeck, normalTail Color
var normalMinLen, normalMaxLen int
var normalMinSpeed, normalMaxSpeed int
var normalSpeedStep int
var normalCharset int

var mutantHead, mutantNeck, mutantTail Color
var mutantMinLen, mutantMaxLen int
var mutantMinSpeed, mutantMaxSpeed int
var mutantSpeedStep int
var mutantCharset int

var frameDuration int
var overlap int
var maxDropsPerColumn int
var reservedHeight int
var syncSpeed int
var cmdPath string

func defaults() {
	// normalHead = Color{255, 153,   0}
	// normalNeck = Color{224,  51,   0}
	// normalTail = Color{22,    5,   5}
	normalHead = Color{136, 204,   0}
	normalNeck = Color{ 51, 153,  51}
	normalTail = Color{  0,  25,   0}
	normalMinSpeed = 2
	normalMaxSpeed = 4
	normalSpeedStep = 4
	normalMinLen = 8
	normalMaxLen = 16
	normalCharset = 2

	mutantHead = Color{204, 153, 255}
	mutantNeck = Color{250, 255, 250}
	mutantTail = Color{  0, 204, 122}
	mutantMinSpeed = 1
	mutantMaxSpeed = 2
	mutantSpeedStep = 1
	mutantMinLen = 12
	mutantMaxLen = 24
	mutantCharset = 0

	frameDuration = 40
	overlap = 5
	maxDropsPerColumn = 50
	reservedHeight = 3
	syncSpeed = 0
	rainStatus = true
}

func readCommandLine() {
	var configPath string
	pflag.StringVarP(&cmdPath, "pipe", "p", "", "[path] pipe for reading commands")
	pflag.StringVarP(&configPath, "config", "c", "", "[path] configuration file")

	pflag.Parse()	

	if configPath != "" {
		readCommandFile(configPath, nil)
	}
}


func readCommandFile(path string, ch chan string) {
	f, err := os.Open(path)
	if err != nil { return }
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if ch == nil {
			processCommand(line)
		} else {
			ch <- line
		}
	}
	close(ch)
}
