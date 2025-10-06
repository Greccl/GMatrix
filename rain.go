package main

import (
	// "math"
	"math/rand/v2"
)





var cols []Column
var syncCount int
var syncAdvance bool
var rainStatus bool


type Cell struct {
	r rune
	value int
}
var scene [][]Cell


func rain_resize() {
	if scrw == 0 { return }
	cols = Reslice(cols, scrw)
	
	for i := range cols {
		cols[i].x = i
		cols[i].resize()
	}
}

func rain_tick() {
	generator_0()
	if syncSpeed > 0 {
		syncCount++
		if syncCount >= syncSpeed {
			syncCount = 0
			syncAdvance = true
		} else {
			syncAdvance = false
		}
	}
	for i := 0; i < len(cols); i++ {
		cols[i].tick()
	}
	scr.Show()
}


/*
var gen1time int

func fgen(a, b float64) int {
	for i:=0; i<3; i++ {
		
	}
}

func f0(x, max int) int {
	sin := math.Sin
	xx := float64(x)
	// y := int(1*sin(1*xx) + 2*sin(0.5*xx) + 0.25*sin(4*xx))
	y := int(2*sin(xx/4))
	if y > max { y = max }
	return y
}

func f1(x, max int) int {
	sin := math.Sin
	xx := float64(x)
	// y := int(1*sin(1*xx) + 2*sin(0.5*xx) + 0.25*sin(4*xx))
	y := int(3*sin(xx/4))+2
	if y > max { y = max }
	return y
}

func generator_1() {
	if gen1time > 0 {
		gen1time--
		return
	}
	
	gen1time = 200
	boxh := 10
	speed := 10
	// var lasty int
	var change bool
	var f = f0
	var lap int
	for x := range cols {
		// y := 0 - rand.IntN(boxh)
		var y int
		if x > scrw/2 { change = true }
		if change {
			switch lap {
				case 0: f = f1
			}
			y = -5 + f(x, boxh)
			change = false
		} else {
			y = -5 + f(x, boxh)
		}
		// lasty = y
		d := cols[x].newDrop()
		d.reset()
		d.makeNormal()
		d.speed = speed
		d.pos = y
	}
}
*/

var genCounter int
var density float32 = 5.0
var spawnLeft float32

func generator_0() {
	genCounter += frameDuration
	if genCounter < 333 { return }
	genCounter = 0

	chunks := float32(scrw) / 10.0 / 3.0
	newCount := chunks + spawnLeft
	newCount *= density
	// newCount *= rand.Float32()

	var mutantSpawned bool
	for ; newCount >= 1.0; newCount -= 1.0 {
		x := rand.IntN((scrw+oddOffset)/2) * 2
		allowDups := rand.Float32() < 0.0001
		if !allowDups {
			initialx := x
			MAX: for max:=0; max<5; max++ {
				x = initialx
				delta := 2
				for {
					if cols[x].count == max { break MAX }
				
					x = x + delta
					// delta *= 2
					// delta *= -2
					if x >= scrw { x = 0 }
					// if x
					if x == initialx { break }
				}
			}
		}

		if cols[x].count >= maxDropsPerColumn {
			continue
		}

		var zero *Drop
		if cols[x].count > 0 {
			zero = &cols[x].drops[cols[x].count-1]
			if zero.pos < overlap { continue }
		}
		
		mut := rand.IntN(1000) < 50
		if mut && mutantSpawned { continue }
		
		d := cols[x].newDrop()
		d.reset()
		if mut {
			d.makeMutant()
			mutantSpawned = true
		} else {
			d.makeNormal()
			if zero != nil {
				d.pos = zero.pos - zero.length
			}
			if d.pos > 0 { d.pos = 0 }
		}
	}

	spawnLeft = newCount
	/*
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
	*/
}
