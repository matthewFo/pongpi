package pong

import (
	"log"
	"os/exec"
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

// Play the given sound
func PlaySound(sound SoundType) {
	cmd := exec.Command("/usr/bin/paplay", string(sound))
	err := cmd.Run()
	if err != nil {
		// log and ignore error since not playing the sound isn't critical
		log.Print(err)
	}
}
