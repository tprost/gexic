package sequencer

import midi "github.com/mattetti/audio/midi"

type Note struct {
	Event *midi.Event
}

func NewNote(event *midi.Event) (*Note, error) {
	note := &Note{}
	note.Event = event
	return note, nil
}
