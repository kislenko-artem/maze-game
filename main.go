package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

type Player struct {
	X float64
	Y float64
	A float64
}

var renderMap =[]string{
	"0000000000000000",
 	"1              1",
	"2              0",
	"3              0",
	"1              0",
	"2              0",
	"3              0",
	"0           1  0",
	"0              0",
	"0      3       0",
	"0              0",
	"0   2          0",
	"0              0",
	"0              0",
	"0              0",
	"0100000000000010",
}
var (
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
	var player = Player{X: 1, Y: 4, A: 0}
	var c float64;
	const fov = 3.14/3.0; // field of view
	colorMap := map[string]color.RGBA{
		"0": yellow,
		"1": green,
		"2": red,
		"3": darkGray,
	}
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{Width: 300, Height: 300, Title: "Game" })
		if err != nil {
			fmt.Println(err)
			return
		}
		defer w.Release()
		for {
			e := w.NextEvent()
			//fmt.Println(e)
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
					}
					if e.Code == key.CodeRightArrow {
						player.A += 0.25
					}
					if e.Code == key.CodeD {
						player.Y += 1
					}
					if e.Code == key.CodeA {
						player.Y -= 1
					}
					if e.Code == key.CodeW {
						player.X += 1
					}
					if e.Code == key.CodeS {
						player.X -= 1
					}
					w.Send(paint.Event{})

				}
			case paint.Event:
				// https://github.com/golang/exp/blob/master/shiny/example/basic/main.go
				colorSign := ""
				op := screen.Src
				w.Fill(image.Rectangle{image.Point{0, 0}, image.Point{300, 300}}, blue1, screen.Src)
				for i := 0; i <= 300; i++ {
					angle := float64(player.A) - fov / 2.0 + fov * float64(i) / float64(300);
					for c = 0.0; c<= 18; c += 0.05 {
						x := player.X + c * math.Sin(angle)
						y := player.Y + c * math.Cos(angle)

						//fmt.Println(c, "x", x, y, string([]byte(renderMap[int(x)])[int(y)]))
						colorSign = string([]byte(renderMap[int(x)])[int(y)])
						if colorSign != " " {
							break
						}
					}
					fmt.Println("angle", angle, "c", c, "player.A", player.A, math.Cos(player.A), math.Sin(player.A))

					if colorSign == " " {
						continue
					}
					size0 := image.Point{1, int(500/(c*2))}
					t0, err := s.NewTexture(size0)
					if err != nil {
						log.Fatal(err)
					}
					t0.Fill(t0.Bounds(), colorMap[colorSign], screen.Src)
					t0Rect := t0.Bounds()

					w.Copy(image.Point{i, int(500/(c*2))}, t0, t0Rect, op, nil)

				}
				w.Publish()
			case error:
				log.Print(e)
			}
		}
	})
}