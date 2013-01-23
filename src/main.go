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

	field = NewGameField(64)
	field.Add(NewSinusoid(field, 1))
	//field.Add(NewHSLWheel(field, 1))

	display := NewWebDisplay(field)
	//display := NewLedDisplay(Settings)

	log.Print("MinFrameTime is ", Settings.MinFrameTime)

	ticks := time.Tick(time.Duration(Settings.MinFrameTime*1000.0) * time.Millisecond)
	for _ = range ticks {

		field.Animate(Settings.MinFrameTime)
		field.RenderTo(display)
	}
}
