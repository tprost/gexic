package sequencer

type Pattern struct {
	Instrument Instrument
	Rows []*Row
}

func NewPattern() (*Pattern, error) {
	pattern := &Pattern{}
	pattern.Rows = []*Row{}
	return pattern, nil
}

func (pattern *Pattern) AddRow(row *Row) {
	pattern.Rows = append(pattern.Rows, row)
}

func (pattern *Pattern) Length() int {
	length := len(pattern.Rows)
	count := 0
	for _, row := range pattern.Rows {
		for _, pattern := range row.Patterns {
			if pattern.Length() + count > length {
				length = pattern.Length() + count
			}
		}
		count++
	}
	return length
}

func (pattern *Pattern) GetRowsAtIndex(index int) []*Row {
	rows := []*Row{}
	if index <= len(pattern.Rows) - 1 {
		rows = append(rows, pattern.Rows[index])
	}
	count := 0
	for _, row := range pattern.Rows {
		for _, pattern := range row.Patterns {
			rows = append(rows, pattern.GetRowsAtIndex(count + index)...)
		}
		count++
	}
	return rows
}
