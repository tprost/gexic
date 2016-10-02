package pattern

import "testing"

func TestNoteValue(t *testing.T) {
	value, _ := NoteValue("c#6")
	if value != 73 {
		t.Error("did not equal 73", value)
	}
}
