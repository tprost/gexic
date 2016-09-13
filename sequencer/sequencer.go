package sequencer

import (
	"github.com/gordonklaus/portaudio"
	midi "github.com/mattetti/audio/midi"
	_ "encoding/binary"
	"fmt"
	"github.com/golang-collections/go-datastructures/queue"
)

var sampleRate int = 44100

var bpm = 100
var position int

var latency int = 50 // ms

type Sequencer struct {
	EventQueue * queue.Queue
	Queue *queue.Queue
	CurrentSample int
	Pattern *Pattern
	Stream *portaudio.Stream
	Tracks []*Track
}

type Track struct {
	CurrentNote *Note
}

func NewTrack() (*Track, error) {
	track := &Track{}
	return track, nil
}

func (track *Track) PlayNote(note *Note) {
	track.CurrentNote = note
	note.Instrument.ProcessEvent(note.Event)
}

type Pattern struct {
	Lines []*Note
}

type Note struct {
	Event *midi.Event
	Instrument Instrument
}

func NewNote(event *midi.Event, instrument Instrument) (*Note, error) {
	note := &Note{}
	note.Instrument = instrument
	note.Event = event
	return note, nil
}

func NewPattern() (*Pattern, error) {
	pattern := &Pattern{}
	pattern.Lines = make([]*Note, 4)
	return pattern, nil
}

func NewSequencer() (*Sequencer, error) {


	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	s := &Sequencer{
	}

	s.Queue = queue.New(100)

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

func (s *Sequencer) AddTrack() *Track {
	track, _ := NewTrack()
	s.Tracks = append(s.Tracks, track)
	return track
}

func (s *Sequencer) ProcessAudio(out []float32) {
	if s.Queue.Empty() {
		for i := range out {
			out[i] = 0
		}
	} else {
		data, _ := s.Queue.Get(int64(len(out)))
		if len(data) == len(out) {
			for i := range out {
				out[i] = data[i].(float32)
			}
		} else {
			for i := range out {
				out[i] = 0
			}
			for i := range data {
				out[i] = data[i].(float32)
			}
		}
	}
}

func (s *Sequencer) QueueSample() {
	var sample float32 = 0
	for _, track := range s.Tracks {
		if track.CurrentNote != nil {
			instrument := track.CurrentNote.Instrument
			trackOut := make([]float32, 1, 1)
			instrument.ProcessAudio(trackOut)
			sample += trackOut[0]
		}
	}
	if sample > 1.0 {
		sample = 1.0
	} else if sample < -1.0 {
		sample = -1.0
	}
	s.Queue.Put(sample)
}
