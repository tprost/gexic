package sequencer

import (
	"github.com/mkb218/gosndfile/sndfile"
	midi "github.com/mattetti/audio/midi"
	"fmt"
	"math"
)

type Instrument interface {
	ProcessAudio(out []float32)
	ProcessEvent(event *midi.Event)
}

type Sampler struct {
	Samples map[uint8][]float32
	Note uint8
	NoteOn bool
	Offset int
	Sample []float32
	FastSample []float32
	SampleRate int
	SampleInfo *sndfile.Info
}

func LinearInterpolation(newSampleRate int, oldSampleRate int, audio []float32) []float32 {
	if newSampleRate == oldSampleRate {
		return audio
	}
	ratio := float64(oldSampleRate)/float64(newSampleRate)
	length := int(float64(len(audio)) * ratio)
	newAudio := make([]float32, length, length)
	for i := range newAudio {
		x := float64(i) / ratio
		r := math.Mod(x, 1)
		d := int(x)

		if d + 1 < len(audio) {
			newAudio[i] = float32(float64(audio[d]) * r + float64(audio[d + 1]) * (1 - r))
		} else {
			newAudio[i] = audio[d]
		}

	}
	return newAudio
}

func NewSampler(filename string) (*Sampler, error) {
	sampler := &Sampler{}
	audio, info, err := LoadSample(filename)
	if err != nil {
		return nil, err
	}
	sampler.Samples = make(map[uint8][]float32)
	sampler.Note = 0
	sampler.NoteOn = false
	sampler.Sample = audio
	sampler.SampleInfo = info
	sampler.SampleRate = 44100
	// for i := 0; i < 100; i++ {
	//	note := i
	//	rate := int(float64(sampler.SampleRate) * math.Pow(2, float64((float64(note) - 60)/12)))
	//	sampler.Samples[sampler.Note] = LinearInterpolation(rate, sampler.SampleRate, sampler.Sample)

	// }
	return sampler, nil
}

func (sampler *Sampler) ProcessAudio(out []float32) {

	if sampler.NoteOn == true {
		sample := sampler.Samples[sampler.Note]
		note := sampler.Note
		if sample == nil {
			rate := int(float64(sampler.SampleRate) * math.Pow(2, float64((float64(note) - 60)/12)))
			sampler.Samples[sampler.Note] = LinearInterpolation(rate, sampler.SampleRate, sampler.Sample)
		}
		sample = sampler.Samples[sampler.Note]
		for i := range out {
			if i + sampler.Offset > len(sample) - 1 {
				out[i] = 0
			} else {
				out[i] = sample[i + sampler.Offset]
			}
		}
	} else {
		for i := range out {
			out[i] = 0
		}
	}
	sampler.Offset += len(out)
}

func (sampler *Sampler) ProcessEvent(event *midi.Event) {
	if (event != nil) {
		if event.MsgType == midi.EventByteMap["NoteOff"] {
			sampler.NoteOn = false
		} else {
			sampler.Note = event.Note
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
