package goadapt

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"syscall"
)

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

func Assert(cond bool, args ...interface{}) {
	var msg string
	if !cond {
		_, file, line, _ := runtime.Caller(1)
		msg = "assertion failed"
		if len(args) == 1 {
			msg += fmt.Sprintf(": %v", args[0])
		}
		if len(args) > 1 {
			msg += ": " + fmt.Sprintf(args[0].(string), args[1:]...)
		}
		e := adaptErr{file, line, msg, nil}
		panic(&e)
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

// convert panic(exitErr) into returned rc and msg
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
