package pong

import (
	"log"
	"os/exec"
	"runtime"
)

// Type of sound identifiers
type SoundType string

// Different sounds
const (
	GAMESTART   SoundType = "sounds/start.wav"
	LEFTBOUNCE            = "sounds/bounce1.wav"
	RIGHTBOUNCE           = "sounds/bounce2.wav"
	MISS                  = "sounds/miss.wav"
	GAMEOVER              = "sounds/gameover.wav"
)

var playWavCommand string

func init() {
	if runtime.GOOS == "windows" {
		playWavCommand = "c:/users/b.green/Desktop/sounder"
	} else {
		playWavCommand = "/usr/bin/paplay"
	}
}

// Play the given sound
func PlaySound(sound SoundType) {
	cmd := exec.Command(playWavCommand, string(sound))
	err := cmd.Start()
	if err != nil {
		// log and ignore error since not playing the sound isn't critical
		log.Print(err)
	}
}
