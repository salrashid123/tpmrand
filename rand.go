package tpmrand

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-tpm/legacy/tpm2"
)

type Reader struct {
	TpmDevice io.ReadWriteCloser
	Scheme    backoff.BackOff
	mu        sync.Mutex
}

func NewTPMRand(conf *Reader) (*Reader, error) {
	if conf.TpmDevice == nil {
		return &Reader{}, fmt.Errorf("Unable to open TPM")
	}

	if conf.Scheme == nil {
		conf.Scheme = backoff.NewExponentialBackOff()
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
	operation := func() (err error) {
		result, err = tpm2.GetRandom(r.TpmDevice, uint16(len(data)))
		if err != nil {
			return err
		}
		copy(data, result)
		return nil
	}

	// dont' know which scheme is better, probably the constant
	//err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	err = backoff.Retry(operation, r.Scheme)
	if err != nil {
		return 0, err
	}

	return len(result), err
}
