package goerr_test

import (
	"fmt"

	"github.com/gtkit/goerr"
)

func ExampleClampMessage() {
	fmt.Println(goerr.ClampMessage("abcdef", 3))
	// Output:
	// abc…
}

func ExampleRedactForLog() {
	fmt.Println(goerr.RedactForLog("call 13900001111"))
	// Output:
	// call [phone]
}
