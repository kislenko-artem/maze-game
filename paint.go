package main

import (
	"github.com/veandco/go-sdl2/img"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"golang.org/x/exp/shiny/screen"
)

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

const (
	fov           = 3.14 / 3.0 // field of view
	wWidth        = 800
	wHeight       = 600
	textureSize   = 64
	offsetTop     = 350
	offsetBottom  = 50
	maxVisibility = 20.0
)

// Player используется для запонимания текущего местоположения игрока
type Player struct {
	PrevX float64 // предыдущие координаты X
	PrevY float64 // предыдущие координаты Y

	X float64 // текущие координаты X
	Y float64 // текущие координаты Y
	A float64 // направление куда игрок смотрит
}

var (
	white  = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

var (
	renderMapCache map[int]map[int]string
	imageData      image.Image
)

func init() {
	var (
		err       error
		imageFile *os.File
	)

	renderMapCache = map[int]map[int]string{}
	// загрузим текстуру
	if imageFile, err = os.Open("assets/textures.png"); err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = imageFile.Close()
	}()
	imageData, _, err = image.Decode(imageFile)
	if err != nil {
		log.Panic("decode: ", err)
	}

	// кеширование  карты
	for x := 0; x < len(renderMap); x++ {
		renderMapCache[x] = map[int]string{}
		for y := 0; y < len(renderMap[x]); y++ {
			renderMapCache[x][y] = string([]byte(renderMap[int(x)])[int(y)])
		}
	}
}

func paintScreen(player *Player, s screen.Screen, w screen.Window) {
	var (
		c float64
	)
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
		return
	}

	size0 := image.Point{wWidth, wHeight}
	img.Init(img.INIT_PNG)
	img.
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

	}

	// отобразим получившееся изображение на экране
	w.Upload(image.Point{0, 0}, imgBuf, imgBuf.Bounds())
	w.Publish()

	// запомним предудещие координаты, чтобы можно было откатиться на них
	player.PrevX = player.X
	player.PrevY = player.Y
}