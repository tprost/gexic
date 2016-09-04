package main

import (
	"github.com/gordonklaus/portaudio"
	_ "encoding/binary"
	"fmt"
	"github.com/mkb218/gosndfile/sndfile"
	"time"
)

var position int

type Sequencer struct {
	Buffer []float32
	Stream	*portaudio.Stream
}

func NewSequencer() (*Sequencer, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	s := &Sequencer{
	}

	stream, err := portaudio.OpenDefaultStream(
		0,
		2,
		float64(44100),
		portaudio.FramesPerBufferUnspecified,
		s.ProcessAudio,
	)

	if err != nil {
		return nil, err
	}

	s.Stream = stream

	return s, nil
}

func (s *Sequencer) Start() {
	fmt.Println("start")
	buffer, _, err := LoadSample("kick.wav")
	s.Buffer = buffer
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	s.Stream.Start()
}

func (s *Sequencer) Close() {
	s.Stream.Close();
}

func (s *Sequencer) ProcessAudio(out []float32) {
	for i := range out {
		var data float32
		if position < len(s.Buffer) {
			data += s.Buffer[position]
			position++
		}
		if data > 1.0 {
			data = 1.0
		}
		out[i] = data

	}
}

// LoadSample loads an audio sample from the passed in filename
// Into memory and returns the buffer.
// Returns an error if there was one in audio processing.
func LoadSample(filename string) ([]float32, *sndfile.Info, error) {
	var info sndfile.Info
	soundFile, err := sndfile.Open(filename, sndfile.Read, &info)
	if err != nil {
		fmt.Printf("Could not open file: %s\n", filename)
		return nil, nil, err
	}

	buffer := make([]float32, 10*info.Samplerate*info.Channels)
	numRead, err := soundFile.ReadItems(buffer)
	if err != nil {
		fmt.Printf("Error reading data from file: %s\n", filename)
		return nil, nil, err
	}

	defer soundFile.Close()

	return buffer[:numRead], &info, nil
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	s, err := NewSequencer()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	s.Start()

		time.Sleep(time.Second)
	s.Close()
}
