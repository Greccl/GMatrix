package main

import (
	// "fmt"
	"math/rand/v2"
	"github.com/Greccl/tcell/v2"
)



type Column struct {
	// state []CellState

	drops []Drop
	count int
	visibleCount int
	x int

	backs Drop
}






func (self *Column) resize() {
	for i := range self.drops {
		self.drops[i].runes = SliceResize(self.drops[i].runes, scrh)
	}
	self.backs.runes = SliceResize(self.backs.runes, scrh)
	// self.state = Reslice(self.state, scrh)
}







func (self *Column) newDrop() *Drop {
	if self.count == len(self.drops) {
		self.drops = append(self.drops, Drop{})
		self.drops[self.count].runes = make([]rune, scrh)
	}
	self.count++
	return &self.drops[self.count-1]
}

func (self *Column) remove(i int) {
	if i < self.count-1 {
		old := self.drops[i].runes
		copy(self.drops[i:], self.drops[i+1:])
		self.drops[self.count-1].runes = old
	}
	self.count--
}

func (self *Column) tick() {
	// self.backTick()
	if self.count == 0 { return }
	self.visibleCount = 0

	var d *Drop

	// Advance
	for i:=self.count-1; i>=0; i-- {
		d = &self.drops[i]

		trueAdvance := false
		if d.mutant || syncSpeed == 0 {
			d.count++
			if d.count >= d.speed { trueAdvance = true }
		} else {
			if syncAdvance { trueAdvance = true }
		}

		if trueAdvance {
			d.pos++
			if d.pos >= 0 { d.dirty = true }
			d.count = 0			
		} else {
			continue
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
	if d.end >= d.pos { d.end = d.pos }
	for i:=self.count-2; i>=0; i-- {
		d = &self.drops[i]
		d.end = d.pos - d.length + 1
		prev := self.drops[i+1].pos + 1
		if d.end < prev { d.end = prev }
		if d.end >= d.pos { d.end = d.pos }
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
		if force && d.pos >= 0 {
			self.visibleCount++
			self.draw(i)
		}
	}

	/*
	s := fmt.Sprintf("%2d", self.count)
	for i := range s {
		scr.SetContent(self.x+i, scrh, rune(s[i]), nil, tcell.StyleDefault)
	}	
	*/
}

func (self *Column) draw(i int) {
	d := &self.drops[i]
	s := tcell.StyleDefault
	var y int

	l := d.pos - d.end + 1
	for p := 0; p <= l; p++ {
		y = d.pos-p
		if y < 0 { return }
		if y >= scrh { continue }
		// can := canDraw(1, self.x, y)
		if p == 0 {
			if d.mutant {
				s.SetForegroundRGB(mutantHead.r, mutantHead.g, mutantHead.b)
			} else {
				if d.lucent {			
					s.SetForegroundRGB(lucentHead.r, lucentHead.g, lucentHead.b)
				} else {
					s.SetForegroundRGB(normalHead.r, normalHead.g, normalHead.b)
				}
			}
		} else if p == l {
			// if can {
				// drawCell(1, self.x, y, ' ', tcell.StyleDefault)
			// }
			releaseCell(1, self.x, y)
			continue
		} else {
			alfa := int32((d.length - p)*1000/d.length)
			var c Color
			if d.mutant {
				c = blend(mutantNeck, mutantTail, 1000-alfa)
			} else {
				if d.lucent {
					
				} else {
					c = blend(normalNeck, normalTail, 1000-alfa)
				}
			}
			s.SetForegroundRGB(c.r, c.g, c.b)
		}
		// if can {
			// scr.SetContent(self.x, y, d.runes[y], nil, s)
			drawCell(1, self.x, y, d.runes[y], s)
		// }
	}
	// if self.x > 0 { cols[self.x-1].backDraw() }
	// if self.x < scrw - 1 { cols[self.x+1].backDraw() }
	
	d.dirty = false
}












var backrunes = []rune{9679, 9670, 9643, 9642, 9702, 9711}

func (self *Column) addBackDrop() {
	d := &self.backs
	d.pos = 0
	d.length = rand.IntN(4) + 4
	d.speed = rand.IntN(3) + 15
	d.back = true
	d.count = 0
	
	for r := range d.runes {
		//d.runes[r] = backrunes[rand.IntN(len(backrunes))]
		d.runes[r] = rand.Int32N(4) + 8756
	}
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

	self.backDraw()
}

func (self *Column) backDraw() {
	d := &self.backs
	if d.length == 0 { return }

	s := tcell.StyleDefault
	var r, g, b int32
	r = backTail.r / 128 //int32(d.length)
	g = backTail.g / 128 //int32(d.length)
	b = backTail.b / 128 //int32(d.length)

	var y int

	// var prev, next *Column
	// if self.x > 0 { prev = &cols[self.x-1] }
	// if self.x < scrw - 1 { next = &cols[self.x+1] }
	
	l := d.pos - d.end + 1
	for p := 0; p <= l; p++ {
		y = d.pos - p
		if y < 0 { return }
		if y >= scrh { continue }
		if p == l {
			scr.SetContent(self.x, y, ' ', nil, tcell.StyleDefault)
			continue
		} else {
			var bright int32 = 510
			// if prev != nil { bright -= prev.state[y].b }
			// if next != nil { bright -= next.state[y].b }
			bright /= 2
			//bright /= 255
			// self.state[y].b = bright
			rr := r * bright
			gg := g * bright
			bb := b * bright
			rr = 41
			gg = 41
			bb = 61
			s.SetForegroundRGB(rr, gg, bb)
		}
		scr.SetContent(self.x, y, d.runes[y], nil, s)
	}
}
