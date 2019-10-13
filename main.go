package adapt

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
)

func Ck(err error) {
	if err != nil {
		panic(err)
	}
}

func Assert(cond bool, args ...interface{}) {
	var txt string
	if !cond {
		_, file, line, _ := runtime.Caller(1)
		if len(args) > 1 {
			txt = fmt.Sprintf(args[0].(string), args[1:]...)
			txt = fmt.Sprintf("assertion failed: %s", txt)
		} else if len(args) > 0 {
			txt = fmt.Sprintf("assertion failed: %v", args[0])
		} else {
			txt = "assertion failed"
		}
		log.Panicf("%s %d: %v", file, line, txt)
	}
}

func Info(txt interface{}, args ...interface{}) {
	_log(log.Printf, txt.(string), args...)
}

func Uerr(txt interface{}, args ...interface{}) {
	_log(log.Panicf, txt.(string), args...)
}

func _log(method func(string, ...interface{}), txt string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	Assert(ok)
	if len(args) > 0 {
		if strings.Contains(txt, "%") {
			txt = fmt.Sprintf(txt, args...)
		} else {
			for _, arg := range args {
				txt += fmt.Sprintf(" %v", arg)
			}
		}
	}
	method("%s %d: %v", file, line, txt)
}

func Passert(t *testing.T, f func()) {
	r := func() {
		if recover() == nil {
			t.Errorf("missing panic")
		}
	}
	defer r()
	f()
}

func Tassert(t *testing.T, cond bool, txt string, args ...interface{}) {
	if !cond {
		debug.PrintStack()
		t.Errorf(txt, args...)
	}
}

func Massert(t *testing.T, str string, substr string) {
	Tassert(t, strings.Contains(str, substr), "can't find %s in %s", substr, str)
}
