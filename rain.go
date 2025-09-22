package main

import (
	"math/rand/v2"
)





var cols []Column



func rain_resize() {
	if scrw == 0 { return }
	cols = Reslice(cols, scrw)
	
	for i := range cols {
		cols[i].x = i
		cols[i].resize()
	}
}

func rain_tick() {
	for i := 0; i < len(cols); i++ {
		cols[i].tick()
	}
	scr.Show()

	newCount := rand.IntN(3)
	for i := 0; i < newCount; i++ {
		x := rand.IntN(scrw/2)
		x *= 2
		allowDups := rand.Float32() < 0.01
		if !allowDups {
			initialx := x
			MAX: for max:=0; max<3; max++ {
				for {
					if cols[x].count == max { break MAX }
					x += 2
					if x >= scrw { x = 0 }
					if x == initialx { break }
				}
			}
		}
		if cols[x].count > 1 {
			continue
		}
		cols[x].addForeDrop()
	}
	
	for i := 0; i < newCount; i++ {
		if rand.IntN(1000) < 25 { // 15
			x := rand.IntN(scrw/2)
			x *= 2
			x++
			if x >= scrw { continue }
			if cols[x].backs.speed > 0 {
				continue
			}
			cols[x].addBackDrop()
		}
	}
}
