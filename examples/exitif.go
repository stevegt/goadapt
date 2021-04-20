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
	rc, msg := run()
	if len(msg) > 0 {
		fmt.Fprint(os.Stderr, msg+"\n")
	}
	os.Exit(rc)
}

func run() (rc int, msg string) {
	defer Halt(&rc, &msg)
	rand.Seed(time.Now().UnixNano())
	err := mid()
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
	r := rand.Intn(4)
	// fmt.Println("r is", r)
	switch r {
	case 0:
		var e syscall.Errno
		e = syscall.EPIPE
		return e
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
