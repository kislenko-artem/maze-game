package main

import (
	"image/color"
	_ "image/png"
	"log"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
)

// Word храним данные мира
type Word struct {
	PrevX       float32
	PastKeyEven time.Time
}

// Player используется для запонимания текущего местоположения игрока
type Player struct {
	PrevX float64 // предыдущие координаты X
	PrevY float64 // предыдущие координаты Y

	X float64 // текущие координаты X
	Y float64 // текущие координаты Y
	A float64 // направление куда игрок смотрит
}

// TODO: выход за текстуру не проверяется
// наша карта. Цифры здесь - это номер картинки из текустуры
var renderMap = []string{
	"0000000000000000",
	"0              0",
	"0              0",
	"0  11111       0",
	"0  1           0",
	"0111           0000000",
	"0              3     0",
	"0              0000000",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0222222222222220",
}

var (
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
	isStop = false
)

const (
	fov           = 3.14 / 3.0 // field of view
	wWidth        = 800
	wHeight       = 600
	textureSize   = 64
	offsetTop     = 350
	offsetBottom  = 50
	maxVisibility = 20.0
)

func main() {
	var (
		player = Player{X: 2, Y: 2, A: 0}
		word   = Word{PrevX: float32(wWidth / 2), PastKeyEven: time.Now()}
	)

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{Width: wWidth, Height: wHeight, Title: "Game"})
		if err != nil {
			return
		}
		defer w.Release()
		go func () {
			for {
				e := w.NextEvent()
				switch e := e.(type) {
				case lifecycle.Event:
					if e.To == lifecycle.StageDead {
						isStop = true
						return
					}
				// что-то вроде управления мышью, хотя из-за ограничений движка работает криво
				// TODO: захват мыши окном
				case mouse.Event:
					mouseEvent(&e, &player, &word)
				// здесь управление игроком
				// TODO: нужно писать свой обработчик нажатия клавиш (чтобы можно было нажимать две одновременно)
				case key.Event:
					keyEvent(&e, &player, &word)
				case error:
					log.Print(e)

				}
			}
		}()

		ticker := time.NewTicker(20 * time.Millisecond)
		timeStart := time.Now()
		var frames float64
		for range ticker.C {
			paintScreen(&player, s, w)
			if isStop == true {
				break
			}
			frames++
			log.Println("fps: ", frames / time.Now().Sub(timeStart).Seconds())
		}
	})
}
