package main

import (
	. "golang-music-stuff/sequencer"
	midi "github.com/mattetti/audio/midi"
	"fmt"
	"time"
)

var sampleRate = 44100


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

	fmt.Println("main")

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

	time.Sleep(time.Second*2)

	player, _ := NewPlayer(s)
	player.Start()


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

	fmt.Println("Ticker stopped")

	s.Close()
}
