package tpmrand

import (
	"crypto/rsa"
	"testing"

	"github.com/google/go-tpm-tools/simulator"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/stretchr/testify/require"
)

func TestRandBytes(t *testing.T) {
	tpmDevice, err := simulator.Get()
	require.NoError(t, err)
	defer tpmDevice.Close()

	tests := []struct {
		name string
		size int
	}{
		{"small", 32},
		{"large", 128},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			randomBytes := make([]byte, tc.size)
			r, err := NewTPMRand(&Reader{
				TpmDevice: tpmDevice,
				//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
			})
			require.NoError(t, err)

			i, err := r.Read(randomBytes)
			require.NoError(t, err)
			require.Equal(t, len(randomBytes), i)
			require.Equal(t, len(randomBytes), tc.size)
		})
	}
}

func TestRandBytesEncrypted(t *testing.T) {
	tpmDevice, err := simulator.Get()
	require.NoError(t, err)
	defer tpmDevice.Close()

	rwr := transport.FromReadWriter(tpmDevice)
	createEKCmd := tpm2.CreatePrimary{
		PrimaryHandle: tpm2.TPMRHEndorsement,
		InPublic:      tpm2.New2B(tpm2.RSAEKTemplate),
	}
	createEKRsp, err := createEKCmd.Execute(rwr)
	require.NoError(t, err)

	randomBytes := make([]byte, 32)
	r, err := NewTPMRand(&Reader{
		TpmDevice:        tpmDevice,
		EncryptionHandle: createEKRsp.ObjectHandle,
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
