package main

import (
	"github.com/mkb218/gosndfile/sndfile"
	"fmt"
)

type Instrument interface {
	ProcessAudio(note *Note, offset int, out []float32)
}

type Sampler struct {
	Sample []float32
	SampleInfo *sndfile.Info
}

func NewSampler(filename string) (*Sampler, error) {
	kick := &Sampler{}
	audio, info, err := LoadSample(filename)
	if err != nil {
		return nil, err
	}
	kick.Sample = audio
	kick.SampleInfo = info
	return kick, nil
}

func (kick *Sampler) ProcessAudio(note *Note, offset int, out []float32) {
	for i := range out {
		if i + offset > len(kick.Sample) - 1 {
			out[i] = 0
		} else {
			out[i] = kick.Sample[i + offset]
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
