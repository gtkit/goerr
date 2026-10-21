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

func ExampleInspect() {
	cause := goerr.New(nil, goerr.StatusParams(), "invalid input")
	input := fmt.Errorf("load profile: %w", cause)
	diagnostic, err := goerr.Inspect(input)
	if err != nil {
		fmt.Println("diagnostic unavailable:", err)
		return
	}
	fmt.Println(diagnostic.Error)
	fmt.Println(diagnostic.Code == goerr.ErrParams, diagnostic.Message)
	fmt.Println(diagnostic.Causes[1].Parent, len(diagnostic.Causes[1].Frames) > 0)
	// Output:
	// load profile: invalid input
	// true invalid input
	// 0 true
}

func ExampleDiagnostic() {
	// 跨包构造值必须使用字段名；通常直接使用 Inspect 的结果。
	diagnostic := goerr.Diagnostic{
		HasError: true,
		Causes:   []goerr.Cause{{Parent: -1}},
	}
	fmt.Println(diagnostic.HasError, diagnostic.Causes[0].Parent)
	// Output:
	// true -1
}
