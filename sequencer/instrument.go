package sequencer

import (
	"github.com/mkb218/gosndfile/sndfile"
	midi "github.com/mattetti/audio/midi"
	"fmt"
	"math"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Instrument interface {
	ProcessAudio(out []float32)
	ProcessEvent(event *midi.Event)
}

type Sampler struct {
	Samples map[uint8][]float32
	Notes map[uint8]bool
	NoteOffsets map[uint8]int
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
	sampler.Notes = make(map[uint8]bool)
	sampler.NoteOffsets = make(map[uint8]int)
	sampler.Sample = audio
	sampler.SampleInfo = info
	sampler.SampleRate = 44100
	return sampler, nil
}

func (sampler *Sampler) ProcessAudio(out []float32) {
	for note, _ := range sampler.Notes {
		sample := sampler.Samples[note]
		if sample == nil {
			rate := int(float64(sampler.SampleRate) * math.Pow(2, float64((float64(note) - 60)/12)))
			sampler.Samples[note] = LinearInterpolation(rate, sampler.SampleRate, sampler.Sample)
		}
		sample = sampler.Samples[note]
		for i := range out {
			offset := sampler.NoteOffsets[note]
			if i + offset <= len(sample) - 1 {
				out[i] += sample[i + offset]
			}
		}
		sampler.NoteOffsets[note] += len(out)
	}
	for i := range out {
		if out[i] > 1.0 {
			out[i] = 1.0
		} else if out[i] < -1.0 {
			out[i] = -1.0
		}
	}


}

func (sampler *Sampler) ProcessEvent(event *midi.Event) {
	if (event != nil) {
		if event.MsgType == midi.EventByteMap["NoteOn"] {
			sampler.Notes[event.Note] = true
			sampler.NoteOffsets[event.Note] = 0
		} else if event.MsgType == midi.EventByteMap["NoteOff"] {
			delete(sampler.Notes, event.Note)
			delete(sampler.NoteOffsets, event.Note)
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



func LoadInstrument(filename string) (Instrument, error) {
	file, _ := ioutil.ReadFile(filename)
	var fileMap map[string]string
	yaml.Unmarshal(file, &fileMap)
	instrument, error := NewSampler(fileMap["sample"])
	return instrument, error
}
