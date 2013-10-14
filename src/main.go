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
		runIntro(buttons, display)
		runOpening(display)
		winner, bounces := runGame(buttons, display)
		go PlayTTS(fmt.Sprint("Game over. Score ", bounces))
		runClosing(buttons, display, winner)
	}
}

// Run an intro animation
func runIntro(buttons *GpioReader, display Display) {

	if runtime.GOOS == "windows" {
		return
	}

	field := NewGameField(Settings.LedCount)
	field.Add(NewSinusoid(field, 1))

	curTime := time.Now()
	prevTime := curTime

	ticks := time.NewTicker(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	defer ticks.Stop()

	for _ = range ticks.C {

		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		if buttons.LeftButton() || buttons.RightButton() {
			return
		}

		field.Animate(dt)
		field.RenderTo(display)
	}
}

// Run an animation to start the game
func runOpening(display Display) {

	field := NewGameField(Settings.LedCount)
	countDown := NewCountdown(field, 2)
	field.Add(countDown)

	go PlaySound(GAMESTART)

	curTime := time.Now()
	prevTime := curTime

	ticks := time.NewTicker(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	defer ticks.Stop()

	for _ = range ticks.C {

		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		field.Animate(dt)

		if countDown.TimeRemaining() <= 0 {
			return
		}

		field.RenderTo(display)
	}
}

// Run the actual game
func runGame(buttons *GpioReader, display Display) (leftPlayerWon bool, totalBounces int) {

	field := NewGameField(64)

	ball := NewBall(field)
	field.Add(ball)

	leftPlayer := NewPlayer(true, Settings.LifeInSeconds, field)
	field.Add(leftPlayer)
	rightPlayer := NewPlayer(false, Settings.LifeInSeconds, field)
	field.Add(rightPlayer)

	curTime := time.Now()
	prevTime := curTime
	totalBounces = 0

	ticks := time.NewTicker(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	defer ticks.Stop()

	for _ = range ticks.C {

		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		leftPlayer.UpdatePaddleActive(buttons.LeftButton())
		rightPlayer.UpdatePaddleActive(buttons.RightButton())
		//leftPlayer.UpdatePaddleActive(true)
		//rightPlayer.UpdatePaddleActive(true)

		field.Animate(dt)

		ball.UpdateOffensiveHide(leftPlayer, rightPlayer)

		playerMissed, bounce := ball.MissedByPlayer(leftPlayer, rightPlayer, Settings.BounceVelocityIncrease)
		if playerMissed != nil {
			ball.ResetPosition(field)
			if playerMissed.DecreaseLife(0.75) {
				return playerMissed == leftPlayer, totalBounces
			}
		}
		if bounce {
			totalBounces++
		}

		field.RenderTo(display)
	}

	panic("Shouldn't get here")
}

// Run an animation showing the winner
func runClosing(buttons *GpioReader, display Display, winner bool) {

	field := NewGameField(Settings.LedCount)
	winnerDisplay := NewWinner(field, winner, 4)
	field.Add(winnerDisplay)

	//go PlaySound(GAMEOVER)

	curTime := time.Now()
	prevTime := curTime

	ticks := time.NewTicker(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	defer ticks.Stop()

	for _ = range ticks.C {

		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		field.Animate(dt)

		if winnerDisplay.TimeRemaining() <= 0 {
			return
		}

		field.RenderTo(display)
	}
}
