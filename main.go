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

var pianoSampler, _ = LoadInstrument("piano.instrument.yaml")

func main() {

	s, err := NewSequencer()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	pattern, _ := NewPattern()
	pattern.Instrument = pianoSampler

	note1, _ := NewNote(midi.NoteOn(0, 60, 50))
	note2, _ := NewNote(midi.NoteOn(0, 59, 50))
	note3, _ := NewNote(midi.NoteOn(0, 58, 50))
	note4, _ := NewNote(midi.NoteOn(0, 57, 50))
	note5, _ := NewNote(midi.NoteOn(0, 56, 50))

	row1, _ := NewRow()
	row2, _ := NewRow()
	row3, _ := NewRow()
	row4, _ := NewRow()
	row5, _ := NewRow()

	row1.AddNote(note1)
	row2.AddNote(note2)
	row3.AddNote(note3)
	row4.AddNote(note4)
	row5.AddNote(note5)

	pattern.AddRow(row1)
	pattern.AddRow(row2)
	pattern.AddRow(row3)
	pattern.AddRow(row4)
	pattern.AddRow(row5)

	s.LoopPattern(pattern)

	player, _ := NewPlayer(s)
	player.Start()

	time.Sleep(time.Second * 5)

	s.Close()
}
