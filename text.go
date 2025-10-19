package main

import (
	"math/rand/v2"
	// "github.com/spf13/pflag"
	"github.com/Greccl/tcell/v2"
)





type TextAnimator interface {
	init(t *Text)
	tick(t *Text)
	draw(t *Text)
}



type Text struct {
	id int
	name string
	runes []rune
	x, y int
	centerx, centery bool
	fg, bg Color
	blends []int32

	anim TextAnimator
	
	animFunc func()
	animName string

	animData_int int
	animData_int32 int32
	animData_bool bool
}

func NewText() *Text {
	t := new(Text)
	t.x = 10
	t.y = len(boxes)
	t.setAnimation("none")
	t.fg = Color{50, 100, 255}
	return t
}

func (t *Text) move(x, y int) {
	t.releaseAll()
	t.x = x
	t.y = y	
}

func (t *Text) setText(str string) {
	str = " " + str + " "
	oldLen := len(t.runes) - 1
	t.runes = []rune(str)
	for oldLen >= len(t.runes) {
		releaseCell(2, t.x+oldLen, t.y)
		oldLen--
	}
	t.blends = make([]int32, len(t.runes))
	t.setAnimation(t.animName)
}

func (t *Text) setAnimation(name string) {
	switch name {
		case "none":
			t.anim = new(TextAnimator0)
			t.anim.init(t)
			// t.animFunc = t.anim_0
			// t.animData_bool = false
		case "basic":
			t.anim = new(TextAnimator1)
			t.anim.init(t)
			// t.animFunc = t.anim_1
			for i := range t.blends {
				t.blends[i] = (rand.Int32N(5)*20) + 100
			}
		case "progresive":
			t.anim = new(TextAnimator2)
			t.anim.init(t)
			t.animData_int = -1
			for i := range t.runes {
				releaseCell(2, t.x+i, t.y)
			}
		default:
			return
	}

	t.animFunc()
	t.animName = name
}

func (t *Text) releaseAll() {
	for x:=0; x<len(t.runes); x++ {
		releaseCell(2, t.x+x, t.y)
	}
}

func (t *Text) autoremove() {
	for i := range boxes {
		if boxes[i] == t {
			boxes = SliceRemove(boxes, i)
			break
		}
	}
	t.releaseAll()
}









func (t *Text) anim_0() {
	if t.animData_bool { return }
	t.animData_bool = true
	s := tcell.StyleDefault
	s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
	s.SetForegroundRGB(t.fg.r, t.fg.g, t.fg.b)
	for i, r := range t.runes {
		drawCell(2, t.x+i, t.y, r, s)
	}
}

func (t *Text) anim_1() {
	s := tcell.StyleDefault
	s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
	for r := range t.blends {
		draw := false
		blender := t.fg
		var alfa int32 = 1000
		if t.blends[r] < 500 {
			t.blends[r] += 10
			blender = Color{}
			alfa = t.blends[r] * 2
			if alfa < 0 { alfa = 0 }
			draw = true
		}
		if t.blends[r] > 500 {
			t.blends[r] -= 25
			blender = Color{255, 200, 200}
			alfa = (1000-t.blends[r]) * 2
			draw = true
		}
		if draw {
			c := blend(t.fg, blender, 1000-alfa)
			s.SetForegroundRGB(c.r, c.g, c.b)
			drawCell(2, t.x+r, t.y, t.runes[r], s)
		}
	}
}


type TextAnimator2 struct {
	
}


func (a *TextAnimator2) init(t *Text) {

}

func (a *TextAnimator2) tick(t *Text) {
	if t.animData_int >= len(t.runes) { return }
	t.animData_int++
	if t.animData_int < len(t.runes) {
		s := tcell.StyleDefault
		s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
		s.SetForegroundRGB(t.fg.r, t.fg.g, t.fg.b)
		drawCell(2, t.x+t.animData_int, t.y, t.runes[t.animData_int], s)
	}
}

func (a *TextAnimator2) draw(t *Text) {

}






