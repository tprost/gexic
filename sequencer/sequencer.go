package sequencer

import (
	"github.com/gordonklaus/portaudio"
	_ "encoding/binary"
	"fmt"
)

var sampleRate int = 44100

var bpm = 44100 // samples per row
var position int

type Sequencer struct {
	Tick int
	Pattern *Pattern
	patternLength int
	Stream *portaudio.Stream
}

func NewSequencer() (*Sequencer, error) {

	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	s := &Sequencer{}

	stream, err := portaudio.OpenDefaultStream(
		0,
		2,
		float64(sampleRate),
		500,
		//		portaudio.FramesPerBufferUnspecified,
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

	return s, nil
}

func (s *Sequencer) Start() {
	s.Stream.Start()
}

func (s *Sequencer) Close() {
	s.Stream.Close();
}

func (s *Sequencer) LoopPattern(pattern *Pattern) {
	s.Pattern = pattern
	s.patternLength = pattern.Length()
}

func (s *Sequencer) ProcessAudio(out Buffer) {

	rows := s.Pattern.GetRowsAtIndex(s.Tick / bpm)
	instrument := s.Pattern.Instrument

	for i := range out {

		// process new position
		if s.Tick > bpm * s.patternLength {
			s.Tick = 0
		}

		if s.Tick % bpm == 0 {
			rows = s.Pattern.GetRowsAtIndex(s.Tick / bpm)
			for _, row := range rows {
				for _, note := range row.Notes {
					instrument.ProcessEvent(note.Event)
				}
			}
		}

		out[i] = 0
		trackOut := make([]float32, 1, 1)
		instrument.ProcessAudio(trackOut)
		out[i] += trackOut[0]

		if out[i] > 1.0 {
			out[i] = 1.0
		} else if out[i] < -1.0 {
			out[i] = -1.0
		}
		s.Tick++
	}
}
