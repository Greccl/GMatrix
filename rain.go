package main

import (
	"math/rand/v2"
	//"github.com/Greccl/tcell/v2"
)


type Rain struct {
	cols []Column
}

func (self *Rain) resize() {
	if scrw == 0 { return }
	if scrw < len(self.cols) {
		self.cols = self.cols[:scrw]
	}
	if scrw > len(self.cols) {
		old := self.cols
		self.cols = make([]Column, scrw)
		copy(self.cols, old)
	}
	
	for i := range self.cols {
		self.cols[i].x = i
		self.cols[i].resize()
	}
}



func (self *Rain) tick() {
	for i := 0; i < len(self.cols); i++ {
		self.cols[i].tick()
	}
	scr.Show()

	newCount := 1
	for i := 0; i < newCount; i++ {
		x := rand.IntN(scrw/2)
		x *= 2
		allowDups := rand.Float32() < 0.01
		if !allowDups {
			initialx := x
			MAX: for max:=0; max<3; max++ {
				for {
					if self.cols[x].count == max { break MAX }
					x += 2
					if x >= scrw { x = 0 }
					if x == initialx { break }
				}
			}
		}
		if self.cols[x].count > 1 {
			continue
		}
		self.cols[x].addForeDrop()
	}
	
	for i := 0; i < newCount; i++ {
		if rand.IntN(1000) < 15 {
			x := rand.IntN(scrw/2)
			x *= 2
			x++
			if x >= scrw { continue }
			if self.cols[x].backs.speed > 0 {
				continue
			}
			self.cols[x].addBackDrop()
		}
	}
}
