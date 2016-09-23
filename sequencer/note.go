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

var NoteMap = map[string]int {
	"a": 58,
	"b": 59,
	"c": 60,
	"d": 61,
	"e": 62,
	"f": 63,
	"g": 64,
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
	}
	if len(str) > 2 {
		i, error := strconv.Atoi(string(str[2]))
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
