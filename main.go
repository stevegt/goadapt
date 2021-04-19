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
	File string
	Line int
	Msg  string
	Err  error
	Rc   int
}

func (e adaptErr) Error() string {
	var s []string
	if len(e.File) > 0 {
		s = append(s, fmt.Sprintf("%s:%d", e.File, e.Line))
	}
	if len(e.Msg) > 0 {
		s = append(s, e.Msg)
	}
	if e.Err != nil {
		s = append(s, fmt.Sprintf("%v", e.Err))
	}
	return strings.Join(s, ": ")
}

func (e adaptErr) Unwrap() error {
	return e.Err
}

func Ck(err error, args ...interface{}) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		msg := formatArgs(args...)
		e := adaptErr{file, line, msg, err, 0}
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
		e := adaptErr{file, line, msg, nil, 0}
		panic(&e)
	}
}

// convert panic into returned err
// see https://github.com/lainio/err2 and https://blog.golang.org/go1.13-errors
func Return(out interface{}, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	// r is an interface{}

	e, ok := r.(*adaptErr)
	if !ok {
		// wasn't us -- let the panic continue
		panic(r)
	}
	// e is an *adaptErr
	// e.Err is the error thrown by lower call

	msg := formatArgs(args...)

	switch res := out.(type) {
	case *error:
		// return a wrapper err
		*res = &adaptErr{Msg: msg, Err: e}
	case *int:
		if e.Rc == 0 {
			// we had an adaptErr panic but no Rc
			log.Println(e)
			*res = 1
		}
		log.Println(e.Msg)
		*res = e.Rc
	}
}

func ExitIf(err, target error, args ...interface{}) {
	// fmt.Printf("%T %T\n", err, target)
	if errors.Is(err, target) {
		rc := int(syscall.EPERM)
		stack := errStack(err)
		// fmt.Printf("%#v\n", stack)
		root := stack[0]
		parent := err
		if len(stack) > 1 {
			parent = stack[1]
		}
		errno, ok := root.(syscall.Errno)
		if ok {
			rc = int(errno)
		}

		msg := formatArgs(args...)
		if len(msg) > 0 {
			msg = fmt.Sprintf("%s: ", msg)
		}
		msg = fmt.Sprintf("%v", parent)

		_, file, line, _ := runtime.Caller(1)
		e := adaptErr{file, line, msg, parent, rc}
		panic(&e)
	}
}

func errStack(e error) (stack []error) {
	stack = append(stack, e)
	child := errors.Unwrap(e)
	if child != nil {
		stack = append(errStack(child), stack...)
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
