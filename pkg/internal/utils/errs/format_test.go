package errs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ExampleFormat_nil() {
	doWork := func(i int) (err error) {
		defer Format(&err, "do work with argument %d", i)
		return nil
	}

	err := doWork(5)

	fmt.Println(err)

	// Output:
	// <nil>
}

func ExampleFormat_error() {
	doWork := func(i int) (err error) {
		defer Format(&err, "do work with argument %d", i)
		return fmt.Errorf("common error")
	}

	err := doWork(5)

	fmt.Println(err)

	// Output:
	// failed to do work with argument 5: common error
}

func ExampleFormatEcho_nil() {
	doWork := func(int) (err error) {
		defer FormatEcho(&err)
		return nil
	}

	err := doWork(5)

	fmt.Println(err)

	// Output:
	// <nil>
}

func ExampleFormatEcho_common_error() {
	doWork := func(i int) (err error) {
		defer FormatEcho(&err)
		return fmt.Errorf("common error for i = %d", i)
	}

	err := doWork(5)

	var echoErr *echo.HTTPError

	fmt.Println(errors.As(err, &echoErr))
	fmt.Println(err)

	// Output:
	// true
	// code=500, message=common error for i = 5
}

func ExampleFormatEcho_errs_error() {
	doWork := func(i int) (err error) {
		defer FormatEcho(&err)
		return Newf(http.StatusBadRequest, "errs.Error for i = %d", i)
	}

	err := doWork(5)

	var echoErr *echo.HTTPError

	fmt.Println(errors.As(err, &echoErr))
	fmt.Println(err)

	// Output:
	// true
	// code=400, message=errs.Error for i = 5
}
