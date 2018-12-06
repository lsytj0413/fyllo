package main

import (
	"fmt"

	"github.com/lsytj0413/ena/algo"
)

func main() {
	s, l := algo.FindNoRepeat("aabcceddabc")
	fmt.Println(s)
	fmt.Println(l)
}
