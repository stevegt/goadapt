package main

import (
	"os"

	. "github.com/stevegt/goadapt"
)

func main() {
	_, err := os.Open("notafile")
	Ck(err)
}
