package waveform

import (
	"errors"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"io"
)

// TODO .................... TODO .................... TODO ....................

type Generator func(reader io.ReadSeeker, duration float64) ([]float64, error)

func NewGenerator() Generator {
	return func(reader io.ReadSeeker, duration float64) ([]float64, error) {
		decoder := wav.NewDecoder(reader)

	}
}
