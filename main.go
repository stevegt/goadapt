package goadapt

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"syscall"
)

// XXX deprecate adaptErr in favor of Wrap and stackTracer from https://pkg.go.dev/github.com/pkg/errors
type adaptErr struct {
	file string
	line int
	msg  string
	err  error
}

func (e adaptErr) Error() string {
	var s []string
	if len(e.file) > 0 {
		s = append(s, fmt.Sprintf("%s:%d", e.file, e.line))
	}
	if len(e.msg) > 0 {
		s = append(s, e.msg)
	}
	if e.err != nil {
		s = append(s, fmt.Sprintf("%v", e.err))
	}
	return strings.Join(s, ": ")
}

// Msg uses UnWrap() to recurse through the err stack, concatenating
// all of the messages found in the stack and returning the result as
// a string.  This function can be used instead of .Error() to get a
// shorter, cleaner message string that doesn't include file and line
// numbers.
func (e adaptErr) Msg() string {
	var parts []string
	msg := e.msg
	if len(msg) > 0 {
		parts = append(parts, msg)
	}
	child := e.Unwrap()
	if child != nil {
		parts = append(parts, errMsg(child))
	}
	return strings.Join(parts, ": ")
}

func (e adaptErr) Unwrap() error {
	return e.err
}

type exitErr struct {
	msg string
	err error
}

func (e exitErr) Error() string {
	var parts []string
	msg := e.msg
	if len(msg) > 0 {
		parts = append(parts, msg)
	}
	child := e.Unwrap()
	if child != nil {
		parts = append(parts, errMsg(child))
	}
	return strings.Join(parts, ": ")
}

func (e exitErr) Unwrap() error {
	return e.err
}

// errRc uses UnWrap() to iterate through the err stack looking for a
// syscall.Errno, and returns that as an int.  Returns 1 if there is
// no syscall.Errno in the stack.
func errRc(err error) (rc int) {
	rc = 1
	e := err
	for {
		e = errors.Unwrap(e)
		if e == nil {
			return
		}
		errno, ok := e.(syscall.Errno)
		if ok {
			rc = int(errno)
			return
		}
	}
}

// errMsg returns a short message describing all errs in stack,
// without filenames and line numbers if possible.
func errMsg(err error) (msg string) {
	switch concrete := err.(type) {
	case *adaptErr:
		msg = concrete.Msg()
	default:
		msg = concrete.Error()
	}
	return
}

/*
// this likely works; leaving this here in case we want it.  in the
// meantime we can just use Ck()
func Raise(i int, args ...interface{}) {
	err := fmt.Errorf("%w: malformed path: %s", syscall.Errno(i), args...)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		msg := formatArgs(args...)
		e := adaptErr{file, line, msg, err}
		panic(&e)
	}
}
*/

func Ck(err error, args ...interface{}) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		msg := formatArgs(args...)
		e := adaptErr{file, line, msg, err}
		panic(&e)
	}
}

func formatArgs(args ...interface{}) (msg string) {
	if len(args) == 1 {
		msg = fmt.Sprintf("%v", args[0])
	}
	if len(args) > 1 {
		msg = fmt.Sprintf(args[0].(string), args[1:]...)
	}
	return
}

/*
func errArgs(args ...interface{}) (err error) {
	var stack []error
	var format string
	for i, arg := args {
		e, ok := arg.(error)
		if ok {
			stack = append(stack, e)
			continue
		}
		// first non-error arg is used as format string
		if len(format) == 0 {
			format = fmt.Sprintf("%v", arg)
			continue
		}
		// remaining args are format values
		e = fmt.Errorf(format, args[i:])
		stack = append(stack, e)
		break
	}
	// wrap in reverse order
	for i := len(stack - 1); i >=0; i-- {
		e = stack[i]
		if err == nil {
			err = e
		} else {
			err = fmt.Errorf("%w", err
		}
	}
}
*/

// Assert takes a bool and zero or more arguments.  If the bool is
// true, then Assert returns.  If the boolean is false, then Assert
// panics.  The panic is of type adaptErr.  The adaptErr contains the
// filename and line number of the caller.  The first argument is used
// as a Sprintf() format string.  Any remaining arguments are provided
// to the Sprintf() as values.  The Sprintf() result is stored as
// adaptErr.msg, to be used later in the adaptErr.Error() string.
func Assert(cond bool, args ...interface{}) {
	if !cond {
		_, file, line, _ := runtime.Caller(1)
		msg := "assertion failed"
		m := formatArgs(args...)
		if len(m) > 0 {
			msg += ": " + m
		}
		err := adaptErr{file, line, msg, nil}
		panic(&err)
	}
}

// ErrnoIf takes a bool and zero or more arguments.  If the bool is
// false, then ErrnoIf returns.  If the boolean is true, then ErrnoIf
// panics.  The panic is of type adaptErr.  The adaptErr contains the
// filename and line number of the caller.  The first argument must be
// of type syscall.Errno; the adaptErr wraps the errno.  The next
// argument is used as a Sprintf() format string.  Any remaining
// arguments are provided to the Sprintf() as values.  The Sprintf()
// result is stored as adaptErr.msg, to be used later in the
// adaptErr.Error() string.
func ErrnoIf(cond bool, errno syscall.Errno, args ...interface{}) {
	if cond {
		_, file, line, _ := runtime.Caller(1)
		msg := formatArgs(args...)
		err := adaptErr{file, line, msg, errno}
		panic(&err)
	}
}

// convert panic into returned err
// see https://github.com/lainio/err2 and https://blog.golang.org/go1.13-errors
func Return(err *error, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	switch concrete := r.(type) {
	case *adaptErr:
		msg := formatArgs(args...)
		*err = &adaptErr{msg: msg, err: concrete}
	case *exitErr:
		msg := formatArgs(args...)
		e := &exitErr{msg: msg, err: concrete}
		panic(e)
	default:
		// wasn't us -- re-raise
		panic(r)
	}
}

// convert panic into returned err on channel
func ReturnChan(errc chan error, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	switch concrete := r.(type) {
	case *adaptErr:
		msg := formatArgs(args...)
		err := adaptErr{msg: msg, err: concrete}
		errc <- err
	case *exitErr:
		msg := formatArgs(args...)
		e := &exitErr{msg: msg, err: concrete}
		panic(e)
	default:
		// wasn't us -- re-raise
		panic(r)
	}
}

// convert panic into returned rc and msg
func Halt(rc *int, msg *string) {
	r := recover()
	if r == nil {
		return
	}
	switch concrete := r.(type) {
	case *adaptErr:
		*rc = errRc(concrete)
		*msg = concrete.Error()
	case *exitErr:
		*rc = errRc(concrete)
		*msg = errMsg(concrete)
	default:
		panic(r)
	}
}

func ExitIf(err, target error, args ...interface{}) {
	if errors.Is(err, target) {
		msg := formatArgs(args...)
		e := &exitErr{msg: msg, err: err}
		panic(e)
	}
}

func errStack(e error) (stack []error) {
	stack = []error{e}
	child := errors.Unwrap(e)
	if child != nil {
		stack = append(stack, errStack(child)...)
	}
	return
}

func Info(msg interface{}, args ...interface{}) {
	_log(log.Printf, msg.(string), args...)
}

func Uerr(msg interface{}, args ...interface{}) {
	_log(log.Panicf, msg.(string), args...)
}

func _log(method func(string, ...interface{}), msg string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	Assert(ok)
	if len(args) > 0 {
		if strings.Contains(msg, "%") {
			msg = fmt.Sprintf(msg, args...)
		} else {
			for _, arg := range args {
				msg += fmt.Sprintf(" %v", arg)
			}
		}
	}
	method("%s %d: %v", file, line, msg)
}
