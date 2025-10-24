package main

import (
	"os"
	"bufio"
	"net"
	"fmt"
	"github.com/spf13/pflag"
)

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





func defaults() {
	// normalHead = Color{255, 153,   0}
	// normalNeck = Color{224,  51,   0}
	// normalTail = Color{22,    5,   5}
	normalHead = Color{136, 204,   0}
	normalNeck = Color{ 51, 153,  51}
	normalTail = Color{  0,  25,   0}
	normalMinSpeed = 2
	normalMaxSpeed = 4
	normalSpeedStep = 3
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
	files      := pflag.StringSliceP("file"    , "f", []string{}, "[path] file for reading commands")
	configPath := pflag.StringP     ("config"  , "c", ""        , "[path] configuration file")
	socket     := pflag.StringP     ("socket"  , "s", "0"       , "create a server socket to read commands")
	pflag.BoolVarP(&kbDriven, "keyboard", "k", false, "enable keyboard events (disables stdin as command source)")
	sockpath   := pflag.Bool("socket-path", false, "print socket path and exit")

	pflag.Parse()	

	// Read config file at the very begining
	if *configPath != "" {
		f, err := os.Open(*configPath)
		if err == nil {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				line := sc.Text()
				parseCommand(line)
			}
		}
	}

	// Read other command files in parallel
	for _, path := range *files {
		go readCommandFile(path)
	}

	// Start the server socket
	if flag := pflag.Lookup("socket"); flag.Changed {
		socketPath := os.TempDir() + "/gmatrix." + *socket + ".sock"
		os.Remove(socketPath)
		if *sockpath {
			fmt.Printf("%s", socketPath)
			os.Exit(0)
		}
		go initServer(socketPath)
	}

	// proccess non-flag arguments as comands
	for _, cmd := range pflag.Args() {
		parseCommand(cmd)
	}
}

func initServer(path string) {
	l, err := net.Listen("unix", path)
	if err != nil { panic(err) }

	defer l.Close()
	defer os.Remove(path)
	for {
		conn, err := l.Accept()
		if err != nil { continue }
		go handleConn(conn)
	}	
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewScanner(conn)

	for reader.Scan() {
		line := reader.Text()
		parseCommand(line)
	}
}

func readCommandFile(path string) {
	f, err := os.Open(path)
	if err != nil { return }
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		parseCommand(line)
	}
}

