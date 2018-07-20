# Errors

A package designed to improve error handling in go.

It revolves around a few simple observations:

- Errors always have a root cause, this is the base case such as io.EOF.
  It allows API's to document specific errors so they can be handled.
- Errors should have a message for humans, this should be informative.
- Errors should optionally have some structured context for debugging.
- Errors should have a stack trace from the origin to where it was handled for debugging
  purposes.


Example:
```
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


```

Prints:

```
A failed: EOF

error trace:

/home/ac/src/errors/cmd/example/main.go:main.A:18 "A failed"
Where:
  "a" = 4
Cause:
/home/ac/src/errors/cmd/example/main.go:main.B:31 "B failed"
Where:
  "b" = 5
Cause:
/home/ac/src/errors/cmd/example/main.go:main.C:41 "C failed"
Where:
  "c" = 6
Cause:
?:? "EOF"
```

The contexts are relatively verbose, but entirely optional.