[PKCS 11 Samples in Go using SoftHSM](https://github.com/salrashid123/go_pkcs11)

### SoftHSM

[softHSM setup](https://github.com/salrashid123/golang-jwt-pkcs11#setup-softhsm)


```bash
export SOFTHSM2_CONF=`pwd`/softhsm.conf
mkdir -p /tmp/soft_hsm/tokens

pkcs11-tool --module /usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so --slot-index=0 --init-token --label="token1" --so-pin="123456"
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so  --label="token1" --init-pin --so-pin "123456" --pin mynewpin
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so --list-token-slots

        Available slots:
        Slot 0 (0x7e66e8ef): SoftHSM slot ID 0x7e66e8ef
        token label        : token1
        token manufacturer : SoftHSM project
        token model        : SoftHSM v2
        token flags        : login required, rng, token initialized, PIN initialized, other flags=0x20
        hardware version   : 2.6
        firmware version   : 2.6
        serial num         : af7ecd99fe66e8ef
        pin min/max        : 4/255
        Slot 1 (0x1): SoftHSM slot ID 0x1
        token state:   uninitialized


pkcs11-tool --module /usr/lib/x86_64-linux-gnu/opensc-pkcs11.so --slot-index=0  --slot=0  --label="token1" --generate-random 50 | xxd -p
```

### TPM 

[PKCS11 Setup for TPM](https://github.com/salrashid123/golang-jwt-pkcs11#tpm)
[golang-jwt for PKCS11](https://github.com/salrashid123/golang-jwt-pkcs11)

```bash
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/libtpm2_pkcs11.so.1 --slot-index=0 --init-token --label="token1" --so-pin="mysopin"
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/libtpm2_pkcs11.so.1 --label="token1" --init-pin --so-pin mysopin --pin mynewpin
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/libtpm2_pkcs11.so.1 --list-token-slots
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/libtpm2_pkcs11.so.1 --label="keylabel1" --login  --pin=mynewpin --id 0  --keypairgen --key-type rsa:2048
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/libtpm2_pkcs11.so.1  --label="keylabel1" --pin mynewpin --generate-random 50 | xxd -p

pkcs11-tool --module /usr/lib/x86_64-linux-gnu/libtpm2_pkcs11.so.1 --list-token-slots
    Available slots:
    Slot 0 (0x1): token1
    token label        : token1
    token manufacturer : GOOG
    token model        : vTPM
    token flags        : login required, rng, token initialized, PIN initialized
    hardware version   : 1.42
    firmware version   : 22.17
    serial num         : 0000000000000000
    pin min/max        : 0/128
    Slot 1 (0x2): 
    token state:   uninitialized
```

### Yubikey


```bash
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/opensc-pkcs11.so --list-token-slots

    Available slots:
        Slot 0 (0x0): Yubico YubiKey OTP+FIDO+CCID 00 00
        token label        : user10@esodemoapp2.com
        token manufacturer : piv_II
        token model        : PKCS#15 emulated
        token flags        : login required, rng, token initialized, PIN initialized
        hardware version   : 0.0
        firmware version   : 0.0
        serial num         : b2955ba170d8b147
        pin min/max        : 4/8

```
[PKCS11 Setup for Yubikey](https://github.com/salrashid123/golang-jwt-pkcs11#setup-yubikey)

[golang-jwt for Yubikey](https://github.com/salrashid123/golang-jwt-yubikey)


---

[Using Go deterministicly generate RSA Private Key with custom io.Reader](https://stackoverflow.com/questions/74869997/using-go-deterministicly-generate-rsa-private-key-with-custom-io-reader) sample [implementaiton](https://go.dev/play/p/1B1KhY_5AHi)