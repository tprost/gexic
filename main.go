package main

import (
	"github.com/gordonklaus/portaudio"
	_ "encoding/binary"
	"fmt"
	"github.com/mkb218/gosndfile/sndfile"
	"time"
)

var SAMPLE_RATE int32

var BPM = 100
var position int

type Sequencer struct {
	CurrentLine int
	Pattern *Pattern
	Stream	*portaudio.Stream
}

type Instrument struct {
	Audio []float32
}

type Pattern struct {
	Lines []*Note
}

type Note struct {
	Instrument *Instrument
}

func NewPattern() (*Pattern, error) {
	pattern := &Pattern{}
	pattern.Lines = make([]*Note, 4)
	return pattern, nil
}

func NewNote() (*Note, error) {
	note := &Note{}
	return note, nil
}

func NewSequencer() (*Sequencer, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	s := &Sequencer{
	}

	kick := &Instrument{}
	audio, info, err := LoadSample("kick.wav")
  SAMPLE_RATE = info.Samplerate
	kick.Audio = audio

	stream, err := portaudio.OpenDefaultStream(
		0,
		2,
		float64(SAMPLE_RATE),
		portaudio.FramesPerBufferUnspecified,
		s.ProcessAudio,
	)

	if err != nil {
		return nil, err
	}

	s.Stream = stream

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	note, _ := NewNote()
	note.Instrument = kick

	pattern, _ := NewPattern()
	pattern.Lines[0] = note
	pattern.Lines[1] = note
	pattern.Lines[2] = note
	pattern.Lines[3] = note

	s.Pattern = pattern

	return s, nil
}

func (s *Sequencer) Start() {
	s.Stream.Start()
}

func (s *Sequencer) Close() {
	s.Stream.Close();
}

func (s *Sequencer) ProcessAudio(out []float32) {
	// fmt.Println("%v", len(out))
	for i := range out {
		var data float32
		if position >= int(SAMPLE_RATE) {
			position = 0
			s.CurrentLine++
			if (s.CurrentLine > 3) {
				s.CurrentLine = 0
			}
		}
		note := s.Pattern.Lines[s.CurrentLine]
		if position < len(note.Instrument.Audio) {
			data += note.Instrument.Audio[position]
		}
		position++
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

	time.Sleep(time.Second * 5)
	s.Close()
}
