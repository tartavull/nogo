package main

import (
    "testing"
)

func TestPackage(t *testing.T) {
    code := "package main\n "
    tokenize(code)

    code = "packagemain\n "
    tokenize(code)

    code = "emain\n "
    _, err := tokenize(code)
    if err != nil {
        t.Error(err)
    }
}
