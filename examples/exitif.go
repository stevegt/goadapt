package main

import (
	"errors"
	"math/rand"
	"os"
	"syscall"
	"time"

	. "github.com/stevegt/goadapt"
)

var FooError = errors.New("foo error")

func main() {
	rand.Seed(time.Now().UnixNano())
	err := mid()
	ExitIf(err, FooError)
	ExitIf(err, syscall.EPIPE, "pipeline %d error", 7)
	ExitIf(err, syscall.ENOENT)
	Ck(err)
}

func adapted() (err error) {
	defer Return(&err)
	err = SomeFunc()
	Ck(err)
	return
}

func mid() (err error) {
	defer Return(&err)
	err = adapted()
	Ck(err)
	return
}

func SomeFunc() (err error) {
	defer Return(&err)
	switch rand.Intn(3) {
	case 0:
		return syscall.EPIPE
	case 1:
		return FooError
	case 2:
		_, err = os.Stat("/notafileordir")
		Ck(err)
	}
	return
}
