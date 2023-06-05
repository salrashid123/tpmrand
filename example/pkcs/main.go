package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/ThalesIgnite/crypto11"
)

func main() {
	randomBytes := make([]byte, 32)

	cctx, err := crypto11.Configure(&crypto11.Config{
		// SoftHSM
		// Pin:        "mynewpin",
		// TokenLabel: "token1",
		// Path:       "/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so",

		// Yubikey
		// Pin:        "123456",
		// TokenLabel: "token1",
		// Path:       "/usr/lib/x86_64-linux-gnu/opensc-pkcs11.so",

		// TPM
		Pin:        "mynewpin",
		TokenLabel: "token1",
		Path:       "/usr/lib/x86_64-linux-gnu/libtpm2_pkcs11.so.1",
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer cctx.Close()

	r, err := cctx.NewRandomReader()
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
