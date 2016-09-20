package main

import (
	. "golang-music-stuff/sequencer"
	"github.com/gordonklaus/portaudio"
	"fmt"
)

type Player struct {
	Sequencer *Sequencer
	AudioBufferer *AudioBufferer
	Stream *portaudio.Stream
}

func NewPlayer(sequencer *Sequencer) (*Player, error) {

	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	p := &Player{}

	p.Sequencer = sequencer
	p.AudioBufferer, _ = NewAudioBufferer(sequencer)

	stream, err := portaudio.OpenDefaultStream(
		0,
		2,
		float64(sampleRate),
		500,
		//		portaudio.FramesPerBufferUnspecified,
		p.AudioBufferer.ProcessAudio,
	)

	if err != nil {
		return nil, err
	}

	p.Stream = stream

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	return p, nil
}

func (p *Player) Start() {
	p.AudioBufferer.Start()
	p.Stream.Start()
}
