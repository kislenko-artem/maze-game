package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
)

// Word храним данные мира
type Word struct {
	PrevX float32
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
		err            error
		imageFile      *os.File
		player         = Player{X: 2, Y: 2, A: 0}
		renderMapCache = map[int]map[int]string{}
		c              float64
	)

	// загрузим текстуру
	if imageFile, err = os.Open("assets/textures.png"); err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = imageFile.Close()
	}()
	imageData, _, err := image.Decode(imageFile)
	if err != nil {
		log.Panic("decode: ", err)
	}
	word := Word{PrevX: float32(wWidth / 2), PastKeyEven: time.Now()}
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{Width: wWidth, Height: wHeight, Title: "Game"})
		if err != nil {
			fmt.Println(err)
			return
		}
		defer w.Release()

		// кеширование  карты
		for x := 0; x < len(renderMap); x++ {
			renderMapCache[x] = map[int]string{}
			for y := 0; y < len(renderMap[x]); y++ {
				renderMapCache[x][y] = string([]byte(renderMap[int(x)])[int(y)])
			}
		}

		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			// что-то вроде управления мышью, хотя из-за ограничений движка работает криво
			// TODO: зхват мышью
			case mouse.Event:
				mouseEvent(&e, &player, &word, w)
			// здесь управление игроком
			// TODO: нужно писать свой обработчик нажатия клавиш (чтобы можно было нажимать две одновременно)
			case key.Event:
				keyEvent(&e, &player, &word, w)
			// здесь рекция на события
			case paint.Event:
			case error:
				log.Print(e)

			}
		}
	})
}
