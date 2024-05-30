package tpmrand

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
)

// Configuration parameters for the TPMRandReader.
type Reader struct {
	TpmDevice        io.ReadWriteCloser // tpm device to use
	Scheme           backoff.BackOff    // backoff retry scheme
	EncryptionHandle tpm2.TPMHandle     // (optional) handle to use for transit encryption
	EncryptionPub    *tpm2.TPMTPublic   // (optional) public key to use for transit encryption
	mu               sync.Mutex
	rwr              transport.TPM
}

// NewTPMRand returns go rand.Reader() from Trusted Platform Module (TPM)
//
//	TPMDevice (io.ReadWriteCloser): The device Handle for the TPM managed by the caller Use either TPMDevice or TPMPath
//	Encrypted (bool): if you want the session encrypted between cpu->tpm
func NewTPMRand(conf *Reader) (*Reader, error) {
	if conf.TpmDevice == nil {
		return &Reader{}, fmt.Errorf("unable to open TPM")
	}
	conf.rwr = transport.FromReadWriter(conf.TpmDevice)
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
		var resp *tpm2.GetRandomResponse
		var sess tpm2.Session
		if r.EncryptionHandle != 0 && r.EncryptionPub != nil {
			sess = tpm2.HMAC(tpm2.TPMAlgSHA256, 16, tpm2.AESEncryption(128, tpm2.EncryptOut), tpm2.Salted(r.EncryptionHandle, *r.EncryptionPub))
		} else {
			sess = tpm2.HMAC(tpm2.TPMAlgSHA256, 16, tpm2.AESEncryption(128, tpm2.EncryptOut))
		}

		resp, err = tpm2.GetRandom{BytesRequested: uint16(len(data))}.Execute(r.rwr, sess)

		if err != nil {
			return err
		}
		result = resp.RandomBytes.Buffer
		copy(data, resp.RandomBytes.Buffer)
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
