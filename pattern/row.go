package pattern

type Row struct {
	Notes []*Note
	Patterns []*Pattern
}

func NewRow() (*Row, error) {
	row := &Row{}
	row.Notes = []*Note{}
	row.Patterns = []*Pattern{}
	return row, nil
}

func (row *Row) AddNote(note *Note) {
	row.Notes = append(row.Notes, note)
}

func (row *Row) AddPattern(pattern *Pattern) {
	row.Patterns = append(row.Patterns, pattern)
}
