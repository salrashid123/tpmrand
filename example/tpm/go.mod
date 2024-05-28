module main

go 1.22

toolchain go1.22.2

require (
	github.com/google/go-tpm v0.9.1-0.20240514145214-58e3e47cd434
	github.com/google/go-tpm-tools v0.3.13-0.20230620182252-4639ecce2aba
	github.com/salrashid123/tpmrand v0.0.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	golang.org/x/sys v0.8.0 // indirect
)

replace github.com/salrashid123/tpmrand => ../../
