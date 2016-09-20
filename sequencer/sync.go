package sequencer

import (
	"github.com/golang-collections/go-datastructures/queue"
	"fmt"
)

type Buffer []float32

type AudioProcessor interface {
	ProcessAudio(out Buffer)
}

type AudioBufferer struct {
	Processor AudioProcessor
	BufferLength int
	Queue *queue.Queue
	Count int
}

func NewAudioBufferer(processor AudioProcessor) (*AudioBufferer, error) {
	bufferer := &AudioBufferer{}
	bufferer.Processor = processor
	bufferer.BufferLength = 10000
	bufferer.Queue = queue.New(int64(bufferer.BufferLength))
	return bufferer, nil
}

func (bufferer *AudioBufferer) ProcessAudio(out Buffer) {

	if bufferer.Queue.Empty() {
		fmt.Println("empty")
		for i := range out {
			out[i] = 0
		}
	} else {
		data, _ := bufferer.Queue.Get(int64(len(out)))
		if len(data) == len(out) {
			fmt.Println("full set")
			for i := range out {
				out[i] = data[i].(float32)
			}
		} else {
			for i := range out {
				out[i] = 0
			}
			for i := range data {
				out[i] = data[i].(float32)
			}
		}
	}
}

func (bufferer *AudioBufferer) Start() {
	for i := 0; i < bufferer.BufferLength; i++ {
		bufferer.Queue.Put(float32(0))
	}
	go func() {
		for {
			if bufferer.Queue.Len() < int64(bufferer.BufferLength) {
				data := make([]float32, 50, 50)
				bufferer.Processor.ProcessAudio(data)
				for _, sample := range data {
					bufferer.Queue.Put(sample)
				}

			}
		}
	}()
}

func (bufferer *AudioBufferer) Stop() {

}
