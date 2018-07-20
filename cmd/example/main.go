package main

import (
	"fmt"
	"io"
	"os"

	"github.com/andrewchambers/errors"
)

func A(a int) error {
	errors := errors.Context([]errors.KV{
		{"a", a},
	})

	err := B(a + 1)
	if err != nil {
		return errors.Wrap(err, "A failed")
	}

	return nil
}

func B(b int) error {
	errors := errors.Context([]errors.KV{
		{"b", b},
	})

	err := C(b + 1)
	if err != nil {
		return errors.Wrap(err, "B failed")
	}
	return nil
}

func C(c int) error {
	errors := errors.Context([]errors.KV{
		{"c", c},
	})

	return errors.Wrap(io.EOF, "C failed")
}

func main() {
	err := A(4)
	fmt.Printf("%s\n\n", err)
	fmt.Printf("error trace:\n\n")
	fmt.Printf("%s", errors.GetTrace(err))
	os.Exit(1)
}
