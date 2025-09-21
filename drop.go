package main

import (
	//"math/rand/v2"
	//"github.com/Greccl/tcell/v2"
)


type Drop struct {
	runes []rune
	mutant bool
	back bool

	length int
	pos int
	head int
	end int
	dirty bool

	speed int
	count int
}




