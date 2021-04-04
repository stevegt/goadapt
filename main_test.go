package adapt

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func passert(t *testing.T, f func()) {
	t.Helper()
	r := func() {
		if recover() == nil {
			t.Errorf("missing panic")
		}
	}
	defer r()
	f()
}

// test boolean condition
func tassert(t *testing.T, cond bool, txt string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Fatalf(txt, args...)
	}
}

func massert(t *testing.T, str string, substr string) {
	t.Helper()
	tassert(t, strings.Contains(str, substr), "can't find %s in %s", substr, str)
}

func TestAssert(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()
	f := func() {
		Assert(false)
	}
	passert(t, f)

	f = func() {
		Assert(false, "foo %s", "bar")
	}
	passert(t, f)

	f = func() {
		Assert(false, "foobar")
	}
	passert(t, f)
}

func TestCheck(t *testing.T) {
	f := func() {
		Ck(errors.New("foo"))
	}
	passert(t, f)

	g := func() (err error) {
		defer Return(&err)
		Ck(errors.New("bar"), "some message")
		return
	}
	err := g()
	massert(t, err.Error(), "some message")

	h := func() (err error) {
		defer Return(&err)
		Ck(errors.New("bar"), "some message with %s", "formatting")
		return
	}
	err = h()
	massert(t, err.Error(), "some message with formatting")

}

func TestReturn(t *testing.T) {

	f := func() (err error) {
		defer Return(&err)
		Ck(errors.New("foo"))
		return
	}
	err := f()
	massert(t, err.Error(), "foo")

	g := func() (err error) {
		defer Return(&err)
		// noop
		return
	}
	err = g()
	tassert(t, err == nil, "err not nil")

	h := func() (err error) {
		defer Return(&err)
		panic("aliens")
	}
	j := func() {
		err := h()
		tassert(t, err == nil, "err not nil")
	}
	passert(t, j)

	k := func() (err error) {
		defer Return(&err, "annotation")
		Ck(errors.New("foo"))
		return
	}
	err = k()
	massert(t, err.Error(), "annotation: ")

	m := func() (err error) {
		defer Return(&err, "annotation with %s", "formatting")
		Ck(errors.New("foo"))
		return
	}
	err = m()
	massert(t, err.Error(), "annotation with formatting: ")

	q := func() (err error) {
		defer Return(&err)
		Ck(errors.New("foo"), "ck annotation")
		return
	}
	err = q()
	massert(t, err.Error(), "ck annotation")

}

type miderr struct{}

func (e miderr) Error() string {
	return "lower error"
}

func (e miderr) Unwrap() error {
	return fmt.Errorf("bottom error")
}

func TestUnwrap(t *testing.T) {

	f := func() (err error) {
		defer Return(&err)
		Ck(&miderr{})
		return
	}
	err := f()
	x := err
loop:
	for i := 0; i < 10; i++ {
		switch y := x.(type) {
		case *AdaptErr:
			// fmt.Printf("%d AdaptErr x %T y %T\n", i, x, y)
			x = y.Unwrap()
		case *miderr:
			// fmt.Printf("%d miderr x %T y %T\n", i, x, y)
			x = y.Unwrap()
		default:
			// fmt.Printf("%d default  x %T y %T\n", i, x, y)
			tassert(t, i == 3, "unwrap depth %d", i)
			tassert(t, x.Error() == "bottom error", "failed bottom unwrap")
			break loop
		}
	}

	var e *AdaptErr
	tassert(t, errors.As(err, &e), "err not unwrapping to AdaptErr")

	var l *miderr
	tassert(t, errors.As(err, &l), "err %T not unwrapping to miderr: %v", err, err)

}

func TestUerr(t *testing.T) {
	f := func() {
		Uerr("foo")
	}
	passert(t, f)
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()

	Info("foo%s", "bar")
	massert(t, buf.String(), "foobar")

	Info("foo", "bar")
	massert(t, buf.String(), "foo bar")

}
