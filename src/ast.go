package main

import (
    "fmt"
    "errors"
)

type Node struct {
}

type File struct {
    Package string
    Imports []string
    Decl []Node
}

func parse(tokens []Token) (File, error) {
    file := File{}
    pos := 0

    fmt.Println(tokens)
    if tokens[pos].id != PACKAGE || tokens[pos+1].id != IDENT {
        return file, errors.New("expected package statement first")
    }
    file.Package = tokens[pos+1].name
    pos = 2
    if tokens[pos].id == IMPORT {
        if tokens[pos+1].id == STRING {
            file.Imports = append(file.Imports, tokens[pos+1].name)
        }
        if tokens[pos+1].id == LPAREN {
            pos += 2
            for tokens[pos].id != RPAREN {
                if tokens[pos].id != STRING {
                    panic("Expected string")
                }
                file.Imports = append(file.Imports, tokens[pos].name)
                pos++
            }
        }
    }

    return file, nil
}

type State struct {
    transition map[Token]State
}

func main () {
    tokens, _ := tokenize(read_file("hello_world.go"))
    file, _ := parse(tokens)
    fmt.Printf("%+v\n", file)
}
