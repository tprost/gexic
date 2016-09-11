package main

import (
	"github.com/gordonklaus/portaudio"
	midi "github.com/mattetti/audio/midi"
	_ "encoding/binary"
	"fmt"
	"time"
)

var SAMPLE_RATE int32 = 44100

var BPM = 100
var position int

type Sequencer struct {
	CurrentSample int
	Pattern *Pattern
	Stream	*portaudio.Stream
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

	stream, err := portaudio.OpenDefaultStream(
		0,
		2,
		float64(SAMPLE_RATE),
		200,
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

	length := len(out)

	for i := range out {
		out[i] = 0
	}

	for _, track := range s.Tracks {
		if track.CurrentNote != nil {
			instrument := track.CurrentNote.Instrument
			trackOut := make([]float32, length, length)
			instrument.ProcessAudio(trackOut)
			for i := range out {
				out[i] += trackOut[i]
			}
		}

	}

	for i := range out {
		if out[i] > 1.0 {
			out[i] = 1.0
		} else if out[i] < -1.0 {
			out[i] = -1.0
		}
	}

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

	track1 := s.AddTrack()
	track2 := s.AddTrack()

	rain, _ := NewSampler("rain.wav")
	kick, _ := NewSampler("kick.wav")

	// c3 := midi.NoteOn(0, 10, 50)
	c4 := midi.NoteOn(0, 20, 50)
	// c3NoteOff := midi.NoteOff(0, 10)
	c4NoteOff := midi.NoteOff(0, 20)

	kickNote, _ := NewNote(c4, kick)
	//	kickNoteOff := NewNote(c4NoteOff, kick)
	rainNote, _ := NewNote(c4, rain)
	rainNoteOff, _ := NewNote(c4NoteOff, rain)

	s.Start()

	track1.PlayNote(kickNote)
	track2.PlayNote(rainNote)

	time.Sleep(time.Second)

	track1.PlayNote(kickNote)

	time.Sleep(time.Second)

	track1.PlayNote(kickNote)

	time.Sleep(time.Second)

	track1.PlayNote(kickNote)
	track2.PlayNote(rainNoteOff)

	time.Sleep(time.Second)

	s.Close()
}
