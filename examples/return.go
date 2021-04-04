package main

import (
	"fmt"

	. "github.com/stevegt/goadapt"
)

func main() {
	err := eg1()
	fmt.Printf("%v\n", err)

	err = eg2()
	fmt.Printf("%v\n", err)

	err = eg3()
	fmt.Printf("%v\n", err)

}

func eg1() (err error) {
	defer Return(&err)
	Assert(false, "foo %s", "bar")
	return
}

func eg2() (err error) {
	defer Return(&err, "some annotation")
	Assert(false, "foo %s", "bar")
	return
}

func eg3() (err error) {
	defer Return(&err, "some annotation with %s", "formatting")
	Assert(false, "foo %s", "bar")
	return
}
