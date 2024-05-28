package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"net"
	"slices"

	//"time"
	"github.com/google/go-tpm-tools/simulator"
	"github.com/google/go-tpm/tpmutil"

	//"github.com/cenkalti/backoff/v4"
	tpmrand "github.com/salrashid123/tpmrand"
)

var (
	tpmPath = flag.String("tpm-path", "/dev/tpm0", "Path to the TPM device (character device or a Unix socket).")
)

var TPMDEVICES = []string{"/dev/tpm0", "/dev/tpmrm0"}

func OpenTPM(path string) (io.ReadWriteCloser, error) {
	if slices.Contains(TPMDEVICES, path) {
		return tpmutil.OpenTPM(path)
	} else if path == "simulator" {
		return simulator.GetWithFixedSeedInsecure(1073741825)
	} else {
		return net.Dial("tcp", path)
	}
}

func main() {

	flag.Parse()

	rwc, err := OpenTPM(*tpmPath)
	if err != nil {
		fmt.Printf("Unable to open TPM at %s", *tpmPath)
	}
	defer rwc.Close()

	randomBytes := make([]byte, 32)
	r, err := tpmrand.NewTPMRand(&tpmrand.Reader{
		TpmDevice: rwc,
		//Encrypted: true,
		//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	// Rand read
	_, err = r.Read(randomBytes)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Random String :%s\n", hex.EncodeToString(randomBytes))

	fmt.Println()

	// // /// RSA keygen
	privkey, err := rsa.GenerateKey(r, 2048)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)
	fmt.Printf("RSA Key: \n%s\n", keyPEM)
}
