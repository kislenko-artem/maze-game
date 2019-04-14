package main

import (
	"golang.org/x/mobile/event/key"
	"log"
	"math"
	"time"

	"golang.org/x/mobile/event/paint"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/mouse"
)

func mouseEvent (e *mouse.Event, player *Player, word *Word, w screen.Window) {
	var direction = "left"
	if word.PrevX > e.X {
		direction = "right"
	}
	word.PrevX = e.X
	// будем контролировать скорость, при зажатой клавиши, т.к. сама библиотека слишком быстро
	// генерирует события
	if time.Now().Sub(word.PastKeyEven) < time.Millisecond * 50 {
		return
	}
	word.PastKeyEven = time.Now()
	if direction == "right" {
		player.A -= 0.25
		if player.A < -6 {
			player.A = 0
		}
	}
	if direction == "left" {
		player.A += 0.25
		if player.A > 6 {
			player.A = 0
		}
	}
	log.Println("mouse",  e.X, player.A, (wWidth / 2), direction)
	w.Send(paint.Event{})

}

func keyEvent (e *key.Event, player *Player, word *Word, w screen.Window) {
	if e.Code == key.CodeEscape {
		return
	}

	// будем контролировать скорость, при зажатой клавиши, т.к. сама библиотека слишком быстро
	// генерирует события
	if time.Now().Sub(word.PastKeyEven) < time.Millisecond*50 {
		return
	}
	word.PastKeyEven = time.Now()

	if e.Direction == key.DirRelease {
		if e.Code == key.CodeLeftArrow {
			player.A -= 0.25
			if player.A < -6 {
				player.A = 0
			}
		}
		if e.Code == key.CodeRightArrow {
			player.A += 0.25
			if player.A > 6 {
				player.A = 0
			}
		}
		cosDirection := math.Round(math.Cos(player.A))
		sinDirection := math.Round(math.Sin(player.A))

		if e.Code == key.CodeD {
			if cosDirection == 0 {
				player.Y -= sinDirection
				player.X -= cosDirection
			} else {
				player.Y += sinDirection
				player.X += cosDirection
			}
		}
		if e.Code == key.CodeA {
			if cosDirection == 0 {
				player.Y += sinDirection
				player.X += cosDirection
			} else {
				player.Y -= sinDirection
				player.X -= cosDirection
			}
		}
		if e.Code == key.CodeW {
			player.Y += cosDirection
			player.X += sinDirection
		}
		if e.Code == key.CodeS {
			player.Y -= cosDirection
			player.X -= sinDirection
		}
		w.Send(paint.Event{})
	}
}
