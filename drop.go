package main

import (
	"math/rand/v2"
	"github.com/Greccl/tcell/v2"
)





type Drop struct {
	runes []rune
	col int
	pos int
	head int
	length int
	speed float32
	count float32
}

func (self *Drop) advance(force bool) bool {
	self.count += self.speed
	if self.count >= 1.0 {
		self.pos++
		if self.pos >= h {
			self.head++
		}
		force = true
		self.count -= 0.99
	}
	if force {
		self.draw()
	}
	if self.head >= len(self.runes) {
		self.speed = 0.0
		return false
	}
	return true
}

func (self *Drop) draw() {
	s := tcell.StyleDefault
	r := tailColor.r / int32(self.length)
	g := tailColor.g / int32(self.length)
	b := tailColor.b / int32(self.length)

	var x, y int

	x = self.col

	for p := 0; p < self.length; p++ {
		y = self.pos-p
		if y < 0 { return }
		if y >= h { continue }
		if p == 0 {
			buffer[x][y]++
			s.SetForegroundRGB(headColor.r, headColor.g, headColor.b)
		} else {
			bright := int32(self.length - p)
			rr := r * bright
			gg := g * bright
			bb := b * bright
			s.SetForegroundRGB(rr, gg, bb)
		}
		scr.SetContent(x, y, self.runes[y], nil, s)
	}
	
	if y >= 0 && y < h {
		buffer[x][y]--
		if buffer[x][y] <= 0 {
			scr.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
		buffer[x][y] = 0
	}
}









type Column struct {
	drops []Drop
	x int
}

func (self *Column) add() {
	var d Drop
	d.speed = (float32(rand.IntN(5)) * 0.5) + 0.5
	d.col = self.x
	d.length = rand.IntN(11) + 10

	i := -1
	if len(self.drops) == 0 {
		self.drops = make([]Drop, 1)
		i = 0
	} else {
		for d := range self.drops {
			if self.drops[d].speed == 0 {
				i = d
				break
			}
		}
	}
	if i == -1 {
		self.drops = append(self.drops, d)
		i = len(self.drops) - 1
	} else {
		self.drops[i] = d
	}
	runes := make([]rune, h)
	for r := range runes {
		c := rand.Int32N(27) + 65
		runes[r] = c
	}
	self.drops[i].runes = runes
}

func (self *Column) tick() {
	force := false
	for i := range self.drops {
		if self.drops[i].speed > 0.0 {
			force = self.drops[i].advance(force) || force
		}
	}
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
	
	newCount := rand.IntN(5)
	if newCount > 2 { newCount-- }
	if newCount > 0 { newCount-- }

	for i := 0; i < newCount; i++ {
		allowDups := rand.IntN(10) < 5
		x := rand.IntN(w/2)
		if !allowDups {
			initialx := x
			for {
				if len(self.cols[x].drops) == 0 { break }
				x++
				if x >= w/2 { x = 0 }
				if x == initialx { break }
			}
		}
		if len(self.cols[x].drops) > 2 {
			continue
		}
		self.cols[x].x = x*2
		self.cols[x].add()
	}
}





