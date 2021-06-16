package main

import (
    "testing"
)

func assert_noerror(t *testing.T, err error) {
    if err != nil {
        t.Error(err)
    }
}

func assert_error(t *testing.T, err error) {
    if err == nil {
        t.Error("expected an error to be returned")
    }
}

func TestPackage(t *testing.T) {
    code := "package main\n "
    _, err := tokenize(code)
    assert_noerror(t, err)

    code = "packagemain\n "
    _, err = tokenize(code)
    assert_error(t, err)

    code = "emain\n "
    _, err = tokenize(code)
    assert_error(t, err)

    code = "package 0main\n "
    _, err = tokenize(code)
    assert_error(t, err)

    code = "package _\n "
    _, err = tokenize(code)
    assert_noerror(t, err)

    code = `//single line comment
package main`
    _, err = tokenize(code)
    assert_noerror(t, err)

    code = `/*some long
multi-line comment 
package main`
    _ , err = tokenize(code)
    assert_error(t, err)

    code = `/*some long
    multi-line comment 

    package main`
    _ , err = tokenize(code)
    assert_error(t, err)

    code = `package main

import "os"

func main() {
    os.Exit(1)
}`

}
