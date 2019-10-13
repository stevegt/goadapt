package adapt

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"
)

func TestAssert(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()
	f := func() {
		Assert(false)
	}
	Passert(t, f)

	f = func() {
		Assert(false, "foo %s", "bar")
	}
	Passert(t, f)
	Massert(t, buf.String(), "foo bar")

	f = func() {
		Assert(false, "foobar")
	}
	Passert(t, f)
	Massert(t, buf.String(), "foobar")
}

func TestCheck(t *testing.T) {
	f := func() {
		Ck(errors.New("foo"))
	}
	Passert(t, f)
}

func TestUerr(t *testing.T) {
	f := func() {
		Uerr("foo")
	}
	Passert(t, f)
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()

	Info("foo%s", "bar")
	Massert(t, buf.String(), "foobar")

	Info("foo", "bar")
	Massert(t, buf.String(), "foo bar")

}
