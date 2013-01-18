package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"
)

// Serve static html page
func htmlPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
	<head><script type="text/javascript"><!--
		function reloadpic()
        {
			document.images["gameBoard"].src = "image/test.png";
			setTimeout(reloadpic, 500);
        }
        setTimeout(reloadpic, 500)
	--></script></head>
	<body><img id="gameBoard" src="image/test.png"/></body>
</html>`)

}

// Return newly generating image
func imageHandler(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-control", "max-age=0, must-revalidate, no-store")

	spacing := 8
	ledCount := 64
	width, height := spacing*ledCount, 32
	image := image.NewRGBA(image.Rect(0, 0, width, height))

	color := color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}

	for ledPos := 0; ledPos < width; ledPos += spacing {
		for y := 0; y < height; y++ {
			image.Set(ledPos, y, color)
		}
	}

	png.Encode(w, image)
	fmt.Println("Generated", r.URL, " in", time.Since(startTime))
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// Application entry point
func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		//defer pprof.StopCPUProfile()
		//defer fmt.Println("EXITING") // this line not executed

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

	http.HandleFunc("/", htmlPageHandler)
	http.HandleFunc("/image/", imageHandler)

	fmt.Println("Server listening on 8080")
	http.ListenAndServe(":8080", nil)
}
