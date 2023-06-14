package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	//"time"

	//"github.com/cenkalti/backoff/v4"
	tpmrand "github.com/salrashid123/tpmrand"
)

func main() {
	randomBytes := make([]byte, 32)
	r, err := tpmrand.NewTPMRand(&tpmrand.Reader{
		TpmDevice: "/dev/tpm0",
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
	fmt.Printf("Random String :%s\n", base64.StdEncoding.EncodeToString(randomBytes))

	fmt.Println()

	// /// RSA keygen
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
