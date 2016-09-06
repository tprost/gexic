package main

import (
	"github.com/mkb218/gosndfile/sndfile"
	midi "github.com/mattetti/audio/midi"
	"fmt"
)

type Instrument interface {
	ProcessAudio(out []float32)
	ProcessEvent(event *midi.Event)
}

type Sampler struct {
	NoteOn bool
	Offset int
	Sample []float32
	SampleInfo *sndfile.Info
}

func NewSampler(filename string) (*Sampler, error) {
	sampler := &Sampler{}
	audio, info, err := LoadSample(filename)
	if err != nil {
		return nil, err
	}
	sampler.NoteOn = false
	sampler.Sample = audio
	sampler.SampleInfo = info
	return sampler, nil
}

func (sampler *Sampler) ProcessAudio(out []float32) {
	for i := range out {
		if sampler.NoteOn == false || i + sampler.Offset > len(sampler.Sample) - 1 {
			out[i] = 0
		} else {
			out[i] = sampler.Sample[i + sampler.Offset]
		}
	}
	sampler.Offset += len(out)
}

func (sampler *Sampler) ProcessEvent(event *midi.Event) {
	if (event != nil) {
		if event.MsgType == midi.EventByteMap["NoteOff"] {
			sampler.NoteOn = false
		} else {
			sampler.NoteOn = true
			sampler.Offset = 0
		}
	}
}

// LoadSample loads an audio sample from the passed in filename
// Into memory and returns the buffer.
// Returns an error if there was one in audio processing.
func LoadSample(filename string) ([]float32, *sndfile.Info, error) {
	var info sndfile.Info
	soundFile, err := sndfile.Open(filename, sndfile.Read, &info)
	if err != nil {
		fmt.Printf("Could not open file: %s\n", filename)
		return nil, nil, err
	}

	buffer := make([]float32, 10*info.Samplerate*info.Channels)
	numRead, err := soundFile.ReadItems(buffer)
	if err != nil {
		fmt.Printf("Error reading data from file: %s\n", filename)
		return nil, nil, err
	}

	defer soundFile.Close()

	return buffer[:numRead], &info, nil
}
