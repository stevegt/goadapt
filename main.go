package adapt

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type AdaptErr struct {
	File string
	Line int
	Msg  string
	Err  error
}

func (e AdaptErr) Error() string {
	var s string
	if len(e.File) > 0 {
		s = fmt.Sprintf("%s:%d: ", e.File, e.Line)
	}
	if len(e.Msg) > 0 {
		s += fmt.Sprintf("%s", e.Msg)
	}
	if e.Err != nil {
		s += fmt.Sprintf(": %v", e.Err)
	}
	return s
}

func (e AdaptErr) Unwrap() error {
	return e.Err
}

func Ck(err error, args ...interface{}) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		var msg string
		if len(args) == 1 {
			msg = fmt.Sprintf("%v", args[0])
		}
		if len(args) > 1 {
			msg = fmt.Sprintf(args[0].(string), args[1:]...)
		}
		e := AdaptErr{file, line, msg, err}
		panic(&e)
	}
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
		e := AdaptErr{file, line, msg, nil}
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
	// r is interface{}

	e, ok := r.(*AdaptErr)
	if !ok {
		// wasn't us -- let the panic continue
		panic(r)
	}
	// e is *AdaptErr
	// e.Err is the original error thrown by lower call

	var msg string
	if len(args) == 1 {
		msg = fmt.Sprintf("%v", args[0])
	}
	if len(args) > 1 {
		msg = fmt.Sprintf(args[0].(string), args[1:]...)
	}

	*err = &AdaptErr{Msg: msg, Err: e}
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
