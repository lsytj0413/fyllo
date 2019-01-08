package main

import (
	"fmt"

	"github.com/lsytj0413/ena/store/fs"
)

func main() {
	store := fs.New()

	w, err := store.Watch("/test", true)
	fmt.Println(err)
	fmt.Println(w)

	r, err := store.Get("/test", false, false)
	fmt.Println(err)
	fmt.Println(r)

	r, err = store.Set("/test", false, "test")
	fmt.Println(err)
	fmt.Println(r)

	r, err = store.Get("/test", false, false)
	fmt.Println(err)
	fmt.Println(r)

	r = <-w.ResultChan()
	fmt.Println(r)

	w.Remove()

	r, err = store.Set("/test", false, "test2")
	fmt.Println(err)
	fmt.Println(r)

	r = <-w.ResultChan()
	fmt.Println(r)

	fmt.Println("store/fs example")
}
