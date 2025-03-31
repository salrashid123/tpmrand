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
	mu               sync.Mutex
	rwr              transport.TPM
	maxDigestBuffer  int
	sess             tpm2.Session
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

	getRsp, err := tpm2.GetCapability{
		Capability:    tpm2.TPMCapTPMProperties,
		Property:      uint32(tpm2.TPMPTMaxDigest),
		PropertyCount: 1,
	}.Execute(conf.rwr) // to do: encrypt
	if err != nil {
		return nil, fmt.Errorf("tpmjwt: failed to run capability %v", err)
	}
	tp, err := getRsp.CapabilityData.Data.TPMProperties()
	if err != nil {
		return nil, fmt.Errorf("tpmjwt: failed to get capability %v", err)
	}
	conf.maxDigestBuffer = int(tp.TPMProperty[0].Value)

	if conf.EncryptionHandle != 0 {
		encryptionPub, err := tpm2.ReadPublic{
			ObjectHandle: conf.EncryptionHandle,
		}.Execute(conf.rwr)
		if err != nil {
			return nil, fmt.Errorf("unable to read Encryption Public key %v", err)
		}
		ePubName, err := encryptionPub.OutPublic.Contents()
		if err != nil {
			return nil, fmt.Errorf("unable to get Encryption Public key contents %v", err)
		}
		conf.sess = tpm2.HMAC(tpm2.TPMAlgSHA256, 16, tpm2.AESEncryption(128, tpm2.EncryptOut), tpm2.Salted(conf.EncryptionHandle, *ePubName))
	} else {
		conf.sess = tpm2.HMAC(tpm2.TPMAlgSHA256, 16, tpm2.AESEncryption(128, tpm2.EncryptOut))
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

	chunkSize := r.maxDigestBuffer

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		operation := func() (err error) {
			var resp *tpm2.GetRandomResponse
			resp, err = tpm2.GetRandom{BytesRequested: uint16(len(chunk))}.Execute(r.rwr, r.sess)
			if err != nil {
				return err
			}
			result = append(result, resp.RandomBytes.Buffer...)
			return nil
		}
		err = backoff.Retry(operation, r.Scheme)
		if err != nil {
			return 0, err
		}

	}
	copy(data, result)
	return len(result), err
}
