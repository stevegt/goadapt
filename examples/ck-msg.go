package main

import (
	"os"

	. "github.com/stevegt/goadapt"
)

func main() {
	fn := "notafile"
	_, err := os.Open(fn)
	Ck(err, "can't find %s", fn)
}
