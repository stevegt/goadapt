package main

import . "github.com/stevegt/goadapt"

func main() {
	Assert(false, "foo %s", "bar")
}
