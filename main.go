package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int32 = 800, 600
var joysticks [16]*sdl.Joystick

func run() int {
	var window *sdl.Window
	var event sdl.Event
	var running bool
	var err error

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	sdl.JoystickEventState(sdl.ENABLE)

	running = true
	timeStart := time.Now()
	var frames float64
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			case *sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyboardEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			case *sdl.JoyAxisEvent:
				fmt.Printf("[%d ms] JoyAxis\ttype:%d\twhich:%c\taxis:%d\tvalue:%d\n",
					t.Timestamp, t.Type, t.Which, t.Axis, t.Value)
			case *sdl.JoyBallEvent:
				fmt.Printf("[%d ms] JoyBall\ttype:%d\twhich:%d\tball:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.Ball, t.XRel, t.YRel)
			case *sdl.JoyButtonEvent:
				fmt.Printf("[%d ms] JoyButton\ttype:%d\twhich:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.Button, t.State)
			case *sdl.JoyHatEvent:
				fmt.Printf("[%d ms] JoyHat\ttype:%d\twhich:%d\that:%d\tvalue:%d\n",
					t.Timestamp, t.Type, t.Which, t.Hat, t.Value)
			default:
				fmt.Printf("Some event\n")
			}
		}

		sdl.Delay(16)
		frames++
		log.Println("fps: ", frames/time.Now().Sub(timeStart).Seconds())
	}

	return 0
}

func main() {
	os.Exit(run())
}