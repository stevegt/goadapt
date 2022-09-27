package main

import (
	"os"

	. "github.com/stevegt/goadapt"
)

func main() {
	err := foo()
	Pl(err)
}

func foo() (err error) {
	defer Return(&err)
	_, err = os.Open("notafile")
	Ck(err)
	return
}
