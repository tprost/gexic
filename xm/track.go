package xm

type Track []*Row

func (track *Track) NewTrack(size int) (*Track, error) {
	track = make([]*Row, size)
}

func (track *Track) RemoveRow() (*Track, error) {

}

func (track *Track) AddRow(row *Row) {

}
