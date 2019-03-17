package main

import (
	"fmt"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"image"
	"image/color"
	"log"
	"math"
)

type Player struct {
	X float64
	Y float64
	A float64
}

var renderMap = []string{
	"0000000000000000",
	"0              1",
	"0              0",
	"0  11111       0",
	"0  1           0",
	"0111           0000000",
	"0                    0",
	"0              0000000",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0              0",
	"0100000000000010",
}
var (
	white    = color.RGBA{0xff, 0xff, 0xff, 0xff}
	blue0    = color.RGBA{0x00, 0x00, 0x1f, 0xff}
	blue1    = color.RGBA{0x00, 0x00, 0x3f, 0xff}
	darkGray = color.RGBA{0x3f, 0x3f, 0x3f, 0xff}
	green    = color.RGBA{0x00, 0x7f, 0x00, 0x7f}
	red      = color.RGBA{0x7f, 0x00, 0x00, 0x7f}
	yellow   = color.RGBA{0x3f, 0x3f, 0x00, 0x3f}

	cos30 = math.Cos(math.Pi / 6)
	sin30 = math.Sin(math.Pi / 6)
)

func main() {
	var (
		player = Player{X: 2, Y: 2, A: 0}
		renderMapCache = map[int]map[int]string{}
		c float64
	)
	const (
		fov = 3.14 / 3.0 // field of view
		wWidth = 800
		wHeight = 600
	)
	colorMap := map[string]color.RGBA{
		"0": yellow,
		"1": green,
		"2": red,
		"3": darkGray,
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


		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

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
					backFroward := &player.X
					rightLeft := &player.Y

					if e.Code == key.CodeD {
						*rightLeft += sinDirection
						*backFroward += cosDirection
					}
					if e.Code == key.CodeA {
						*rightLeft -= sinDirection
						*backFroward -= cosDirection
					}
					if e.Code == key.CodeW {
						*rightLeft += cosDirection
						*backFroward += sinDirection
					}
					if e.Code == key.CodeS {
						*rightLeft -= cosDirection
						*backFroward -= sinDirection
					}
					//fmt.Println("a", player.A, "cos", cosDirection, "sin", sinDirection, "x", player.X, "y", player.Y)
					w.Send(paint.Event{})

				}
			case paint.Event:
				colorSign := ""
				if player.X < 1 {
					player.X = 1
				}
				if player.Y < 1 {
					player.Y = 1
				}
				if player.X > float64(len(renderMap) - 2) {
					player.X = float64(len(renderMap) - 2)
				}
				if player.Y > float64(len(renderMap[int(player.X)]) - 2) {
					player.Y = float64(len(renderMap[int(player.X)]) - 2)
				}

				//if renderMapCache[int(player.X) + 1][int(player.Y)] != " " {
				//	player.X = player.X - 1
				//}
				//
				//if renderMapCache[int(player.X) - 1][int(player.Y)] != " " {
				//	player.X = player.X + 1
				//}
				//
				//if renderMapCache[int(player.X)][int(player.Y) - 1] != " " {
				//	player.Y = player.Y + 1
				//}
				//
				//if renderMapCache[int(player.X)][int(player.Y) + 1] != " " {
				//	player.Y = player.Y - 1
				//}
				size0 := image.Point{wWidth, wHeight, }
				imgBuf, err := s.NewBuffer(size0)
				if err != nil {
					log.Fatal(err)
				}
				defer imgBuf.Release()
				img := imgBuf.RGBA()

				for i := 0; i <= wWidth; i+=1 {
					angle := float64(player.A) - fov/2.0 + fov*float64(i)/float64(wWidth)
					for c = 0.0; c <= 20; c += 0.01 {
						x := player.X + c*math.Sin(angle)
						y := player.Y + c*math.Cos(angle)
						colorSign = renderMapCache[int(x)][int(y)]
						if colorSign != " " {
							break
						}

					}
					if colorSign == " " || colorSign == "" {
						continue
					}
					sizeY := int(wHeight / (c * math.Cos(angle - player.A))) + (wHeight / 5)
					if sizeY > wHeight {
						sizeY = wHeight
					}
					if sizeY < 0 {
						sizeY = 0
					}

					for b := 0; b < wHeight; b++ {
						if b > sizeY {
							img.SetRGBA(i, b, white)
							continue
						}
						img.SetRGBA(i, b, colorMap[colorSign])
					}

				}
				w.Upload(image.Point{0, 0}, imgBuf, imgBuf.Bounds())
				w.Publish()
			case error:
				log.Print(e)
			}
		}
	})
}
