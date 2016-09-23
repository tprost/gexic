package main

import (
	. "golang-music-stuff/sequencer"
	"fmt"
	"time"
)

var sampleRate = 44100

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

var pianoSampler, _ = LoadInstrument("piano.instrument.yaml")

func main() {

	s, err := NewSequencer()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	pattern, _ := LoadPattern("test.pattern")
	pattern.Instrument = pianoSampler

	s.LoopPattern(pattern)

	player, _ := NewPlayer(s)
	player.Start()

	time.Sleep(time.Second * 10)

	s.Close()
}
