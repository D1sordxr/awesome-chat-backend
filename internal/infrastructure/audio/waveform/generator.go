package waveform

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/go-audio/wav"
)

const (
	DefaultBuckets = 64
	maxBuckets     = 4096
)

var (
	ErrInvalidWAV  = errors.New("waveform: not a valid wav stream")
	ErrEmptyStream = errors.New("waveform: stream contains no pcm samples")
)

type Generator func(reader io.ReadSeeker, buckets int) ([]float64, error)

func NewGenerator() Generator {
	return func(reader io.ReadSeeker, buckets int) ([]float64, error) {
		const op = "waveform.Generator"

		switch {
		case buckets <= 0:
			buckets = DefaultBuckets
		case buckets > maxBuckets:
			buckets = maxBuckets
		}

		decoder := wav.NewDecoder(reader)
		if !decoder.IsValidFile() {
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidWAV)
		}

		if err := decoder.Rewind(); err != nil {
			return nil, fmt.Errorf("%s: rewind: %w", op, err)
		}

		buf, err := decoder.FullPCMBuffer()
		if err != nil {
			return nil, fmt.Errorf("%s: read pcm: %w", op, err)
		}
		if buf == nil || len(buf.Data) == 0 {
			return nil, fmt.Errorf("%s: %w", op, ErrEmptyStream)
		}

		return peaks(buf.Data, buf.SourceBitDepth, buckets), nil
	}
}

func peaks(samples []int, bitDepth, buckets int) []float64 {
	if buckets > len(samples) {
		buckets = len(samples)
	}

	full := fullScale(bitDepth)
	out := make([]float64, buckets)

	for i := range out {
		start := i * len(samples) / buckets
		end := (i + 1) * len(samples) / buckets
		if end <= start {
			end = start + 1
		}

		var peak int
		for _, s := range samples[start:end] {
			if s < 0 {
				s = -s
			}
			if s > peak {
				peak = s
			}
		}

		out[i] = math.Min(float64(peak)/full, 1)
	}

	return out
}

func fullScale(bitDepth int) float64 {
	if bitDepth <= 0 || bitDepth > 64 {
		bitDepth = 16
	}
	return math.Ldexp(1, bitDepth-1)
}
