package main

import (
    "fmt"
    "io/ioutil"
    "errors"
)

type state uint

// All possible states for the tokenizer
const (
    Initial state = iota
    CommentSingleLine
    CommentMultiLine
    Package
    Space
)

//All tokens
const (
)

type State struct {
    StateType state
    AllowedTransition []state
}

func NotImplemented() {
    panic("Not Implemented")
}

func Failure() {
}

func peek(pattern string, corpse string, pos int) bool {
    if pos + len(pattern) > len(corpse) {
        return false
    }
    return corpse[pos:pos+len(pattern)] == pattern
}

func tokenize(code string) ([]string, error) {
    state := Initial
    pos := 0
    col := 0
    row := 0
    var tokens []string
    for pos < len(code) {
        if code[pos] == '\n' {
            col=0
            row++
        }

        if state == Initial {
            if peek("package", code, pos) {
                pos += len("package")
                tokens = append(tokens, "package")
                state = Space
            } else if peek("//", code, pos) {
                pos += len("//")
                state = CommentSingleLine
            } else if peek("/*", code, pos) {
                pos += len("/*")
                state = CommentMultiLine
            } else {
                return nil, errors.New(fmt.Sprintf("%d:%d Expected 'package' or comment.",row, col))
            }
        }

        pos += 1
        col += 1
    }
    return tokens, nil
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func read_file(filename string)  string {
    full_path := "../test_programs/" + filename
    content, err := ioutil.ReadFile(full_path)
    check(err)
    return string(content)
}

func main () {
    tokenize(read_file("hello_world.go"))
}
