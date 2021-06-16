package main

import (
    "fmt"
    "io/ioutil"
    "errors"
    "regexp"
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

func match(pattern string, corpse string, pos int)  []string {
    re := regexp.MustCompile(pattern)
    return re.FindStringSubmatch(corpse[pos:])
}

func tokenize(code string) ([]string, error) {
    pos := 0
    col := 0
    row := 0
    var tokens []string
    for pos < len(code) {
        if code[pos] == '\n' {
            col=0
            row++
        }
        fmt.Println("reminder:", string(code[pos:]))
        if m := match(`^(package)\s+([_a-zA-Z][_a-zA-Z0-9]*)\s*`, code, pos); len(m) > 0 {
            // Matches package statement
            pos += len(m[0])
            tokens = append(tokens, m[1:]...)
        } else if m := match(`^\/\/[^\n]*\n\s*`, code, pos); len(m) > 0 {
            // Matches single line comments
            pos += len(m[0])
        } else if m := match(`^\/\*[^\*\/]*(\*\/)?\s*`, code, pos) ; len(m) > 0 {
            // Matches multi-line comments and detects if not closed.
            if len(m[1]) == 0 {
                return nil, errors.New(fmt.Sprintf("%d:%d Multi-line comment never closes.",row, col))
            }
            pos += len(m[0])
        } else {
            return nil, errors.New(fmt.Sprintf("%d:%d Expected 'package' or comment.",row, col))
        }
        col += pos
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
