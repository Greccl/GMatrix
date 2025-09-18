package main

import (
	"math/rand/v2"
	"github.com/Greccl/tcell/v2"
)





type Drop struct {
	runes []rune

	length int
	pos int
	head int
	end int
	dirty bool

	speed int
	count int
}





type Column struct {
	drops []Drop
	count int
	x int
}

func (self *Column) add() {
	var d *Drop

	if self.count == len(self.drops) {
		self.drops = append(self.drops, Drop{})
		self.drops[self.count].runes = make([]rune, h)
	}
	self.count++
	d = &self.drops[self.count-1]

	d.pos = 0
	d.length = rand.IntN(16) + 10
	d.speed = rand.IntN(4) + 1
	d.count = 0
	for r := range d.runes {
		switch face {
			case 0:
				d.runes[r] = rand.Int32N(27) + 65
			case 1:
				d.runes[r] = rand.Int32N(2) + 48
			case 2:
				d.runes[r] = rand.Int32N(93) + 33
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

func (self *Column) tick() {
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
	if self.drops[0].end >= h {
		self.remove(0)
		scr.SetContent(self.x, h-1, ' ', nil, tcell.StyleDefault)
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
	r := tailColor.r / int32(d.length)
	g := tailColor.g / int32(d.length)
	b := tailColor.b / int32(d.length)

	var y int

	l := d.pos - d.end + 1
	for p := 0; p <= l; p++ {
		y = d.pos-p
		if y < 0 { return }
		if y >= h { continue }
		if p == 0 {
			s.SetForegroundRGB(headColor.r, headColor.g, headColor.b)
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









type Rain struct {
	cols []Column
}

func (self *Rain) resize(w, h int) {
	w /= 2
	if w == 0 { w = 1 }
	if w < len(self.cols) {
		self.cols = self.cols[:w]
	}
	if w > len(self.cols) {
		old := self.cols
		self.cols = make([]Column, w)
		copy(self.cols, old)
	}
}



func (self *Rain) tick() {
	for i := 0; i < len(self.cols); i++ {
		self.cols[i].tick()
	}
	scr.Show()

	newCount := 1

	for i := 0; i < newCount; i++ {
		x := rand.IntN(w/2)
		allowDups := rand.Float32() < 0.01
		if !allowDups {
			initialx := x
			MAX: for max:=0; max<3; max++ {
				for {
					if self.cols[x].count == max { break MAX }
					x++
					if x >= w/2 { x = 0 }
					if x == initialx { break }
				}
			}
		}
		// x = 0
		if self.cols[x].count > 2 {
			continue
		}
		self.cols[x].x = x*2
		self.cols[x].add()
	}
}





