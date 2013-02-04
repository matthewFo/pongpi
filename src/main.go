package main

import (
	"flag"
	"fmt"
	"log"
	_ "log"
	_ "math"
	"os"
	"os/signal"
	. "pong"
	. "pong/draw"
	"runtime"
	"runtime/pprof"
	"time"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
var webDisplay = flag.Bool("webdisplay", false, "use webhost on localhost:8080 for the display")

// Application entry point
func main() {

	Settings.Read()

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		//defer pprof.StopCPUProfile() // line not executed because of how exiting from http.ListenAnServe is a process kill

		// capture ctrl+c and stop CPU profiler
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				fmt.Printf("captured %v, stopping profiler and exiting..", sig)
				pprof.StopCPUProfile()
				os.Exit(1)
			}
		}()

		fmt.Println("Start profiling")
	}

	log.Print("MinFrameTime is ", Settings.MinFrameTime)

	var display Display
	if *webDisplay || runtime.GOOS == "windows" {
		display = NewWebDisplay(Settings)
	} else {
		display = NewLedDisplay(Settings)
	}

	buttons := NewGpioReader(Settings)

	// should intro and game(play / dead / win) be different states in statemachine?

	// loop forever
	for {
		//runIntro(buttons, display)
		runGame(buttons, display)
	}
}

// Run an intro animation
func runIntro(buttons *GpioReader, display Display) {

	field := NewGameField(Settings.LedCount)
	field.Add(NewSinusoid(field, 1))

	curTime := time.Now()
	prevTime := curTime

	ticks := time.Tick(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	for _ = range ticks {

		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		if buttons.LeftButton() || buttons.RightButton() {
			return
		}

		field.Animate(dt)
		field.RenderTo(display)
	}
}

// Run the actual game
func runGame(buttons *GpioReader, display Display) {

	field := NewGameField(64)
	//field.Add(NewStepFunction(32.0, 32.0, RGBA{100, 0, 0, 255}, 1))

	ball := NewBall(field)
	field.Add(ball)

	leftPlayer := NewPlayer(true, field)
	field.Add(leftPlayer)
	rightPlayer := NewPlayer(false, field)
	field.Add(rightPlayer)

	curTime := time.Now()
	prevTime := curTime

	ticks := time.Tick(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	for _ = range ticks {

		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		// game play (normal ball is moving gameplay)
		// death (a player has just scored)
		//		showing current score
		//		explosion animation
		//		after 1 second cool down, losing player presses key to go back to playing state
		//  should death be a substate of intro, intro could just be it's own

		// for intro, want seperate field, with different animations going
		// game play and death can share a field

		leftPlayer.UpdatePaddleActive(buttons.LeftButton())
		rightPlayer.UpdatePaddleActive(buttons.RightButton())
		//leftPlayer.UpdatePaddleActive(true)
		//rightPlayer.UpdatePaddleActive(true)

		field.Animate(dt)

		playerMissed := ball.MissedByPlayer(leftPlayer, rightPlayer)
		if playerMissed != nil {
			ball.ResetPosition(field)
		}

		field.RenderTo(display)
	}
}
