package main

import (
	"github.com/gordonklaus/portaudio"
	midi "github.com/mattetti/audio/midi"
	_ "encoding/binary"
	"fmt"
	"time"
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

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

var pianoSampler, _ = NewSampler("note.wav")

func playPianoNote(track *Track, note int) {
	noteOn := midi.NoteOn(0, note, 50)
	pianoNote, _ := NewNote(noteOn, pianoSampler)
	track.PlayNote(pianoNote)
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

	c5 := midi.NoteOn(0, 60, 50)
	c5NoteOff := midi.NoteOff(0, 60)

	rainNote, _ := NewNote(c5, rain)
	rainNoteOff, _ := NewNote(c5NoteOff, rain)

	track1.PlayNote(rainNote)

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			// fmt.Println("Tick at", t)
			t = t
			for i := 0; i < sampleRate; i++ {
				s.QueueSample()
			}
			// fmt.Println("queue is of size", s.Queue.Len())
		}
	}()

	time.Sleep(time.Second*2)

	s.Start()

	playPianoNote(track2, 60)
	time.Sleep(time.Millisecond*250)
	playPianoNote(track2, 62)
	time.Sleep(time.Millisecond*250)
	playPianoNote(track2, 63)
	time.Sleep(time.Millisecond*250)
	playPianoNote(track2, 64)
	time.Sleep(time.Millisecond*250)
	playPianoNote(track2, 57)
	time.Sleep(time.Millisecond*250)
	playPianoNote(track2, 55)
	time.Sleep(time.Millisecond*250)
	playPianoNote(track2, 52)
	time.Sleep(time.Second * 5)
	track1.PlayNote(rainNoteOff)
	ticker.Stop()
	fmt.Println("Ticker stopped")

	s.Close()
}
