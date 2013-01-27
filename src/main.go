package main

import (
	"flag"
	"fmt"
	"log"
	_ "log"
	_ "math"
	"os"
	"os/signal"
	"pong"
	. "pong"
	. "pong/draw"
	"runtime/pprof"
	"time"
)

// Global variables
var (
	field *pong.GameField
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// Application entry point
func main() {

	Settings.Read()

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
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

	field = NewGameField(64)
	field.Add(NewSinusoid(field, 1))
	//field.Add(NewHSLWheel(field, 1))

	field.Add(NewBall(field))

	leftPlayer := NewPlayer(true, field)
	field.Add(leftPlayer)
	rightPlayer := NewPlayer(false, field)
	field.Add(rightPlayer)

	display := NewWebDisplay(field)
	//display := NewLedDisplay(field, Settings)

	buttons := NewGpioReader(Settings)

	log.Print("MinFrameTime is ", Settings.MinFrameTime)

	curTime := time.Now()
	prevTime := curTime

	ticks := time.Tick(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	for _ = range ticks {

		curTime = time.Now()
		dt := curTime.Sub(prevTime).Seconds()
		prevTime = curTime

		leftPlayer.UpdateVisible(buttons.LeftButton())
		rightPlayer.UpdateVisible(buttons.RightButton())

		field.Animate(dt)
		field.RenderTo(display)
	}
}
