package tpmrand

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/google/go-tpm/tpm2"
)

type Reader struct {
	TpmDevice string
	mu        sync.Mutex
	rwc       io.ReadWriteCloser
}

func NewTPMRand(conf *Reader) (*Reader, error) {
	var err error
	conf.rwc, err = tpm2.OpenTPM(conf.TpmDevice)
	if err != nil {
		return &Reader{}, fmt.Errorf("tpm-rand: Public: Unable to Open TPM: %v", err)
	}
	return conf, nil
}

func (r *Reader) Read(data []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(data) > math.MaxUint16 {
		return 0, errors.New("tpm-rand: Number of bytes to read exceeds cannot math.MaxInt16")
	}

	var result []byte
	result, err = tpm2.GetRandom(r.rwc, uint16(len(data)))
	if err != nil {
		return 0, err
	}
	copy(data, result)
	return len(result), err
}

func (r *Reader) Shutdown() error {
	return r.rwc.Close()
}
