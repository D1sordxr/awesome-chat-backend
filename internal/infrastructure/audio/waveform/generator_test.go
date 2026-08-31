package waveform

import (
	"bytes"
	"errors"
	"io"
	"math"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

const (
	testSampleRate = 8000
	testBitDepth   = 16
	testChannels   = 1
)

type memWriteSeeker struct {
	buf []byte
	pos int64
}

func (m *memWriteSeeker) Write(p []byte) (int, error) {
	end := m.pos + int64(len(p))
	if end > int64(len(m.buf)) {
		grown := make([]byte, end)
		copy(grown, m.buf)
		m.buf = grown
	}
	copy(m.buf[m.pos:end], p)
	m.pos = end
	return len(p), nil
}

func (m *memWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var next int64
	switch whence {
	case io.SeekStart:
		next = offset
	case io.SeekCurrent:
		next = m.pos + offset
	case io.SeekEnd:
		next = int64(len(m.buf)) + offset
	default:
		return 0, errors.New("bad whence")
	}
	if next < 0 {
		return 0, errors.New("negative position")
	}
	m.pos = next
	return next, nil
}

func encodeWAV(t *testing.T, samples []int) io.ReadSeeker {
	t.Helper()

	out := &memWriteSeeker{}
	enc := wav.NewEncoder(out, testSampleRate, testBitDepth, testChannels, 1)

	buf := &audio.IntBuffer{
		Format:         &audio.Format{NumChannels: testChannels, SampleRate: testSampleRate},
		Data:           samples,
		SourceBitDepth: testBitDepth,
	}
	if err := enc.Write(buf); err != nil {
		t.Fatalf("encode wav: %v", err)
	}
	if err := enc.Close(); err != nil {
		t.Fatalf("close encoder: %v", err)
	}

	return bytes.NewReader(out.buf)
}

func TestGeneratorRisingRamp(t *testing.T) {
	const total = 1024

	samples := make([]int, total)
	for i := range samples {
		samples[i] = i * math.MaxInt16 / total
	}

	got, err := NewGenerator()(encodeWAV(t, samples), 8)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(got) != 8 {
		t.Fatalf("bucket count = %d, want 8", len(got))
	}

	for i := 1; i < len(got); i++ {
		if got[i] <= got[i-1] {
			t.Errorf("bucket %d = %f, want greater than previous %f", i, got[i], got[i-1])
		}
	}
	if got[len(got)-1] > 1 {
		t.Errorf("last bucket = %f, want <= 1", got[len(got)-1])
	}
}

func TestGeneratorSilenceIsZero(t *testing.T) {
	got, err := NewGenerator()(encodeWAV(t, make([]int, 512)), 4)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	for i, v := range got {
		if v != 0 {
			t.Errorf("bucket %d = %f, want 0", i, v)
		}
	}
}

func TestGeneratorNormalizesFullScalePeak(t *testing.T) {
	samples := make([]int, 256)
	samples[100] = math.MaxInt16

	got, err := NewGenerator()(encodeWAV(t, samples), 1)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("bucket count = %d, want 1", len(got))
	}
	if math.Abs(got[0]-1) > 0.001 {
		t.Errorf("peak bucket = %f, want ~1", got[0])
	}
}

func TestGeneratorDefaultsBucketCount(t *testing.T) {
	got, err := NewGenerator()(encodeWAV(t, make([]int, 4096)), 0)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(got) != DefaultBuckets {
		t.Errorf("bucket count = %d, want %d", len(got), DefaultBuckets)
	}
}

func TestGeneratorClampsBucketsToSampleCount(t *testing.T) {
	got, err := NewGenerator()(encodeWAV(t, make([]int, 10)), 64)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(got) != 10 {
		t.Errorf("bucket count = %d, want 10", len(got))
	}
}

func TestGeneratorRejectsNonWAV(t *testing.T) {
	_, err := NewGenerator()(bytes.NewReader([]byte("this is definitely not a wav file")), 8)
	if !errors.Is(err, ErrInvalidWAV) {
		t.Errorf("err = %v, want ErrInvalidWAV", err)
	}
}
