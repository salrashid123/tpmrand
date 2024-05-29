package tpmrand

import (
	"crypto/rsa"
	"testing"

	"github.com/google/go-tpm-tools/simulator"
	"github.com/stretchr/testify/require"
)

func TestRandBytes(t *testing.T) {
	tpmDevice, err := simulator.Get()
	require.NoError(t, err)
	defer tpmDevice.Close()

	randomBytes := make([]byte, 32)
	r, err := NewTPMRand(&Reader{
		TpmDevice: tpmDevice,
		//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
	})
	require.NoError(t, err)

	i, err := r.Read(randomBytes)
	require.NoError(t, err)
	require.Equal(t, len(randomBytes), i)
	require.Equal(t, len(randomBytes), 32)
}

func TestRandBytesEncrypted(t *testing.T) {
	tpmDevice, err := simulator.Get()
	require.NoError(t, err)
	defer tpmDevice.Close()

	randomBytes := make([]byte, 32)
	r, err := NewTPMRand(&Reader{
		TpmDevice: tpmDevice,
		Encrypted: true,
		//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
	})
	require.NoError(t, err)

	_, err = r.Read(randomBytes)
	require.NoError(t, err)

	require.Equal(t, len(randomBytes), 32)
}

func TestRSAKey(t *testing.T) {
	tpmDevice, err := simulator.Get()
	require.NoError(t, err)
	defer tpmDevice.Close()

	r, err := NewTPMRand(&Reader{
		TpmDevice: tpmDevice,
		//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
	})
	require.NoError(t, err)

	// /// RSA keygen
	privkey, err := rsa.GenerateKey(r, 2048)
	require.NoError(t, err)

	rsaPubKey := privkey.PublicKey
	require.Equal(t, 2048, rsaPubKey.Size()*8)
}
