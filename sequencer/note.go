package sequencer

import midi "github.com/mattetti/audio/midi"
import "strings"
import "strconv"
import "errors"

type Note struct {
	Event *midi.Event
}

func NewNote(event *midi.Event) (*Note, error) {
	note := &Note{}
	note.Event = event
	return note, nil
}

func (note *Note) IsNoteOff() bool {
	return note.Event != nil && note.Event.MsgType == midi.EventByteMap["NoteOff"]
}

func (note *Note) IsNoteOn() bool {
	return note.Event != nil && note.Event.MsgType == midi.EventByteMap["NoteOn"]
}

func (note *Note) ToNoteOff() (*Note, error) {
	return NewNote(midi.NoteOff(0, int(note.Event.Note)))
}

var NoteMap = map[string]int {
	"a": 69,
	"b": 71,
	"c": 60,
	"d": 62,
	"e": 64,
	"f": 65,
	"g": 67,
}

func NoteValue(str string) (int, error) {
	if str == "" {
		return 0, errors.New("empty string")
	}
	value, exists := NoteMap[strings.ToLower(string(str[0]))]
	if !exists {
		return 0, errors.New("could not determine note value")
	}
	if len(str) > 1 {
		char2 := string(str[1])
		if char2 == "#" || char2 == "s" {
			value++
		}
		if char2 == "b" || char2 == "f" {
			value--
		}
		var num string
		if len(str) > 2 {
			num = string(str[2])
		} else {
			num = string(str[1])
		}
		i, error := strconv.Atoi(num)
		if error != nil || i < 0 || i > 9 {
			return 0, errors.New("third character was something unexpected")
		}
		value = value + 12 * (i-5)
	}
	return value, nil
}

func NewNoteFromString(str string) (*Note, error) {
	value, error := NoteValue(str)
	if error != nil {
		return nil, error
	}
	note, error := NewNote(midi.NoteOn(0, value, 50))
	if error != nil {
		return nil, error
	}
	return note, nil
}

func NewNoteOffFromString(str string) (*Note, error) {
	value, error := NoteValue(str)
	if error != nil {
		return nil, error
	}
	note, error := NewNote(midi.NoteOff(0, value))
	if error != nil {
		return nil, error
	}
	return note, nil
}
