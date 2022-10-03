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
	"0              0",
	"0              0000000",
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

func main() {
	var (
		err            error
		imageFile      *os.File
		player         = Player{X: 2, Y: 2, A: 0}
		renderMapCache = map[int]map[int]string{}
		c              float64
	)
	const (
		fov           = 3.14 / 3.0 // field of view
		wWidth        = 800
		wHeight       = 600
		textureSize   = 64
		offsetTop     = 350
		offsetBottom  = 50
		maxVisibility = 20.0
		step          = 0.25
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

		var prevX = float32(wWidth / 2)
		pastKeyEven := time.Now()
		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			// что-то вроде управления мышью, хотя из-за ограничений движка работает криво
			case mouse.Event:
				var direction = "left"
				if prevX > e.X {
					direction = "right"
				}
				prevX = e.X
				// будем контролировать скорость, при зажатой клавиши, т.к. сама библиотека слишком быстро
				// генерирует события
				if time.Now().Sub(pastKeyEven) < time.Millisecond*50 {
					break
				}
				pastKeyEven = time.Now()
				if direction == "right" {
					player.A -= step
					if player.A < -6 {
						player.A = 0
					}
				}
				if direction == "left" {
					player.A += step
					if player.A > 6 {
						player.A = 0
					}
				}
				log.Println("mouse", e.X, player.A, (wWidth / 2), direction)
				w.Send(paint.Event{})
			// здесь управление игроком
			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

				// будем контролировать скорость, при зажатой клавиши, т.к. сама библиотека слишком быстро
				// генерирует события
				if time.Now().Sub(pastKeyEven) < time.Millisecond*50 {
					break
				}
				pastKeyEven = time.Now()

				if e.Direction == key.DirRelease {
					if e.Code == key.CodeLeftArrow {
						player.A -= step
						if player.A < -6 {
							player.A = 0
						}
					}
					if e.Code == key.CodeRightArrow {
						player.A += step
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
			// здесь рекция на события
			case paint.Event:
				mapSign := ""
				if player.X < 1 {
					player.X = 1
				}
				if player.Y < 1 {
					player.Y = 1
				}
				if player.X > float64(len(renderMap)-2) {
					player.X = float64(len(renderMap) - 2)
				}
				if player.Y > float64(len(renderMap[int(player.X)])-2) {
					player.Y = float64(len(renderMap[int(player.X)]) - 2)
				}

				if renderMapCache[int(player.X)][int(player.Y)] != " " {
					player.X = player.PrevX
					player.Y = player.PrevY
					break
				}

				size0 := image.Point{wWidth, wHeight}
				imgBuf, err := s.NewBuffer(size0)
				if err != nil {
					log.Fatal(err)
				}
				defer imgBuf.Release()
				img := imgBuf.RGBA()

				for i := 0; i <= wWidth; i += 1 {

					// вычислим угол, под которым смотрим на мир
					angle := float64(player.A) - fov/2.0 + fov*float64(i)/float64(wWidth)
					var xWall, yWall float64

					// на расстоянии видимости вычислим символ карты, на которую попадаем под этим углом
					for c = 0.0; c <= maxVisibility; c += 0.01 {
						xWall = player.X + c*math.Sin(angle)
						yWall = player.Y + c*math.Cos(angle)
						mapSign = renderMapCache[int(xWall)][int(yWall)]
						if mapSign != " " {
							break
						}

					}
					if mapSign == " " || mapSign == "" {
						continue
					}

					// определим длину текущей линии
					sizeY := int(wHeight/(c*math.Cos(angle-player.A))) + (wHeight / 5)
					if sizeY > wHeight {
						sizeY = wHeight
					}
					if sizeY < 0 {
						sizeY = 0
					}

					go func(sizeY, wHeight, i int, xWall, yWall float64, mapSign string) {
						for b := 0; b < wHeight; b++ {
							if b < wHeight-(sizeY+offsetTop) {
								img.SetRGBA(i, b, white)
								continue
							}
							if b > (sizeY - offsetBottom) {
								img.SetRGBA(i, b, white)
								continue
							}

							// нужен, чтобы выбрать правильное изображение из текстуры
							koef, _ := strconv.Atoi(mapSign)
							koef = koef * textureSize

							// здесь соотнесем текущие размеры и размеры текстуры
							yPic := int(b * textureSize / (sizeY - offsetBottom))
							xPic := int((xWall - float64(int(xWall))) * textureSize)
							if xPic == 0 || xPic == (textureSize-1) {
								xPic = int((yWall - float64(int(yWall))) * textureSize)
							} else {
								yPic = int(b-(wHeight-(sizeY+offsetTop))) * textureSize / (sizeY - offsetBottom - (wHeight - (sizeY + offsetTop)))
							}

							// нарисуем пиксель на изображении
							colorR, colorG, colorB, colorA := imageData.At(xPic+koef, yPic).RGBA()
							img.SetRGBA(i, b, color.RGBA{uint8(colorR), uint8(colorG), uint8(colorB), uint8(colorA)})
						}
					}(sizeY, wHeight, i, xWall, yWall, mapSign)

				}

				// отобразим получившееся изображение на экране
				w.Upload(image.Point{0, 0}, imgBuf, imgBuf.Bounds())
				w.Publish()

				// запомним предудещие координаты, чтобы можно было откатиться на них
				player.PrevX = player.X
				player.PrevY = player.Y
			case error:
				log.Print(e)

			}
		}
	})
}
