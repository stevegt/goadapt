package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"

	. "github.com/stevegt/goadapt"
)

var FooError = errors.New("foo error")

func main() {
	rc, err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
	}
	os.Exit(rc)
}

func run() (rc int, err error) {
	defer Return(&rc)
	defer Return(&err)
	rand.Seed(time.Now().UnixNano())
	err = mid()
	ExitIf(err, FooError)
	ExitIf(err, syscall.EPIPE, "pipeline %d error", 7)
	ExitIf(err, syscall.ENOENT)
	Ck(err)
	return
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
	switch rand.Intn(4) {
	case 0:
		return syscall.EPIPE
	case 1:
		return FooError
	case 2:
		_, err = os.Stat("/notafileordir")
		Ck(err)
	case 3:
		Assert(false, "lksadjfslkjf dsalkjf")
	}
	return
}
