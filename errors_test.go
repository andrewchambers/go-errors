package errors

import (
	"fmt"
	"io"
)

func ExampleWrapped() {
	err := Wrap(io.EOF, "corrupt file")
	fmt.Println(err.Error())
	if RootCause(err) == io.EOF {
		fmt.Println("root cause was indeed io.EOF")
	}

	// Output: corrupt file: EOF
	// root cause was indeed io.EOF
}

func ExampleContext() {
	errors := Context([]KV{
		{"id", 5},
	})

	// This error now displays the context values in its stack traces.
	_ = errors.Wrap(io.EOF, "corrupt file")
}

func ExampleGetTrace() {
	err := Wrap(io.EOF, "initial error")

	errors := Context([]KV{
		{"id", 5},
	})

	err = errors.Wrap(err, "another error")

	errors = Context([]KV{
		{"another", "foobar"},
	})

	err = errors.Wrapf(err, "another %s", "error")
	fmt.Println(GetTrace(err).String())
}
