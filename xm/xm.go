package xm

type word int16
type dword int32

type Row struct {
	Note byte
	Instrument byte
	Volume byte
	EffectType byte
	EffectParameter byte
}


type Instrument struct {
	Name string
	Type byte // no idea what this is
	Samples []*Sample
	SampleMap [96]byte
	VolumeEnvelope [48]byte
	PanningEnvelope [48]byte
	NumberOfVolumePoints byte
	NumberOfPanningPoints byte
	VolumeSustainPoint byte
	VolumeLoopStartPoint byte
	VolumeLoopEndPoint byte
	PanningSustainPoint byte
	PanningLoopStartPoint byte
	PanningLoopEndPoint byte
	VolumeType byte
	PanningType byte
	VibratoType byte
	VibratoSweep byte
	VibratoDepth byte
	VibratoRate byte
	VolumeFadeout word
	Reserved word
}


type Sample struct {
	Length dword
	LoopStart dword
	LoopLength dword
	Volume byte
	Finetune byte
	Type byte
	Panning byte
	RelativeNoteNumber byte
	Reserved byte
	Name string
	Data []float32
}

type Module struct {

}
