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

	log.Print("Creating field")
	field = pong.NewGameField(64)
	//field.Add(pong.NewSinusoid(math.Pi*4, pong.RGBA{255, 0, 0, 255}, 1))
	field.Add(pong.NewHSLWheel(field, 1))

	log.Print("Creating display")
	//display := pong.NewWebDisplay(field)
	display := pong.NewLedDisplay(Settings)

	prevTime := time.Now()
	curTime := time.Now()

	log.Print("MinFrameTime is ", Settings.MinFrameTime)

	ticks := time.Tick(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	for _ = range ticks {

		prevTime = curTime
		curTime = time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		//log.Print("Loop", dt)

		field.Animate(dt)
		field.RenderTo(display)

		// limit to a max FPS
		// if dt < Settings.MinFrameTime {
		// 	sleepFor := Settings.MinFrameTime - dt
		// 	sleepForDuration := time.Duration(sleepFor*1000.0) * time.Millisecond
		// 	//log.Print("Sleep ", sleepFor, sleepForDuration)
		// 	time.Sleep(time.Duration(sleepFor*1000.0) * time.Millisecond)
		// }
	}
}

// loop the game
func gameLoop() {

}
