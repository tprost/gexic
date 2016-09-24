package sequencer

import midi "github.com/mattetti/audio/midi"
import "strings"
import "strconv"
import "errors"
import "fmt"

type Note struct {
	Event *midi.Event
}

func NewNote(event *midi.Event) (*Note, error) {
	note := &Note{}
	note.Event = event
	return note, nil
}

var NoteMap = map[string]int {
	"a": 57,
	"b": 59,
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
	fmt.Println(value)
	if error != nil {
		return nil, error
	}
	return note, nil
}
