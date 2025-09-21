package main

import (
	"math/rand/v2"
	"github.com/Greccl/tcell/v2"
)

type Column struct {
	drops []Drop
	count int
	x int

	backs Drop
}


func (self *Column) resize() {
	for i := range self.drops {
		self.drops[i].runes = Reslice(self.drops[i].runes, scrh)
	}
	self.backs.runes = Reslice(self.backs.runes, scrh)
	self.count = 0
}


var backrunes = []rune{9679, 9670, 9643, 9642, 9702, 9711}

func (self *Column) addBackDrop() {
	d := &self.backs
	d.pos = 0
	d.length = rand.IntN(10) + 10
	d.speed = rand.IntN(5) + 10
	d.back = true
	d.count = 0
	
	for r := range d.runes {
		//d.runes[r] = backrunes[rand.IntN(len(backrunes))]
		d.runes[r] = rand.Int32N(4) + 8756
	}
}

func (self *Column) addForeDrop() {
	var d *Drop

	if self.count == len(self.drops) {
		self.drops = append(self.drops, Drop{})
		self.drops[self.count].runes = make([]rune, scrh)
	}
	self.count++
	d = &self.drops[self.count-1]

	d.pos = 0
	d.mutant = false
	d.length = rand.IntN(12) + 6
	d.speed = rand.IntN(5) + 5
	if rand.IntN(100) < 10 { // make it mutant
		d.speed = rand.IntN(2) + 1
		d.mutant = true
		d.length += 10
	}
	d.count = 0
	for r := range d.runes {
		switch face {
			case 0:
				d.runes[r] = rand.Int32N(27) + 65
			case 1:
				d.runes[r] = rand.Int32N(2) + 48
			case 2:
				d.runes[r] = rand.Int32N(93) + 33
			case 3:
				d.runes[r] = rand.Int32N(96) + 0x30A0
			case 4:
				d.runes[r] = rand.Int32N(96) + 0x3040
			case 5:
				d.runes[r] = rand.Int32N(144) + 0x0370
		}
	}
}

func (self *Column) remove(i int) {
	if i < self.count-1 {
		old := self.drops[i].runes
		copy(self.drops[i:], self.drops[i+1:])
		self.drops[self.count-1].runes = old
	}
	self.count--
}

func (self *Column) backTick() {
	d := &self.backs
	
	if d.speed == 0 { return }
	
	d.count++
	if d.count >= d.speed {
		d.pos++
		d.count = 0
	} else {
		return
	}

	// Ending
	d.end = d.pos - d.length + 1
	if d.end < 0 { d.end = 0 }
	if d.end >= scrh {
		d.speed = 0
	}

	// Drawing
	s := tcell.StyleDefault
	var r, g, b int32
	r = backTail.r // int32(d.length)
	g = backTail.g // int32(d.length)
	b = backTail.b // int32(d.length)

	var y int

	l := d.pos - d.end + 1
	for p := 0; p <= l; p++ {
		y = d.pos - p
		if y < 0 { return }
		if y >= scrh { continue }
		if p == l {
			scr.SetContent(self.x, y, ' ', nil, tcell.StyleDefault)
			continue
		} else {
			/*
			bright := int32(d.length - p)
			if bright < 4 { bright = 4}
			if bright > int32(d.length) { bright = int32(d.length) }
			rr := r * bright
			gg := g * bright
			bb := b * bright
			s.SetForegroundRGB(rr, gg, bb)
			*/
			s.SetForegroundRGB(r, g, b)
		}
		scr.SetContent(self.x, y, d.runes[y], nil, s)
	}
}

func (self *Column) tick() {
	self.backTick()
	if self.count == 0 { return }

	var d *Drop

	// Advance
	for i:=0; i<self.count; i++ {
		d = &self.drops[i]
		d.count++
		if d.count >= d.speed {
			d.pos++
			d.dirty = true
			d.count = 0
		}
		// Check overlaps
		if i > 0 {
			if d.pos >= self.drops[i-1].pos {
				self.remove(i-1)
				i--
			}
		}
	}

	// Endings
	d = &self.drops[self.count-1]
	d.end = d.pos - d.length + 1
	if d.end < 0 { d.end = 0 }
	for i:=self.count-2; i>=0; i-- {
		d = &self.drops[i]
		d.end = d.pos - d.length + 1
		prev := self.drops[i+1].pos + 1
		if d.end < prev { d.end = prev }
	}

	// Remove completed
	if self.drops[0].end >= scrh {
		self.remove(0)
		scr.SetContent(self.x, scrh-1, ' ', nil, tcell.StyleDefault)
	}

	// Drawing
	force := false
	for i:=0; i<self.count; i++ {
		force = force || self.drops[i].dirty
		if force { self.draw(i) }
	}
}

func (self *Column) draw(i int) {
	d := &self.drops[i]
	s := tcell.StyleDefault
	var r, g, b int32
	if d.mutant {
		r = mutantTail.r / int32(d.length)
		g = mutantTail.g / int32(d.length)
		b = mutantTail.b / int32(d.length)
	} else {
		r = tailColor.r / int32(d.length)
		g = tailColor.g / int32(d.length)
		b = tailColor.b / int32(d.length)
	}

	var y int

	l := d.pos - d.end + 1
	for p := 0; p <= l; p++ {
		y = d.pos-p
		if y < 0 { return }
		if y >= scrh { continue }
		if p == 0 {
			if d.mutant {
				s.SetForegroundRGB(mutantHead.r, mutantHead.g, mutantHead.b)
			} else {
				s.SetForegroundRGB(headColor.r, headColor.g, headColor.b)
			}
		} else if p == l {
			scr.SetContent(self.x, y, ' ', nil, tcell.StyleDefault)
			continue
		} else {
			bright := int32(d.length - p)
			rr := r * bright
			gg := g * bright
			bb := b * bright
			s.SetForegroundRGB(rr, gg, bb)
		}
		scr.SetContent(self.x, y, d.runes[y], nil, s)
	}
	d.dirty = false
}