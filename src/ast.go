package main

import (
    "fmt"
    "errors"
)

type Node interface {
    Type() NodeType
    Children() []Node
}

type NodeType uint

const (
	FuncType NodeType = iota
    ValueSpec
    TypeSpec
)

type NodeFunc struct {
    name string
    arg_name []string
    arg_type []string
    ret_type []string
    body []Node
}

func (n NodeFunc) Type() NodeType {
    return FuncType
}

func (n NodeFunc) Children() []Node {
    return n.body
}

type NodeTypeSpec struct {
}

func (n NodeTypeSpec) Type() NodeType {
    return TypeSpec
}

func (n NodeTypeSpec) Children() []Node {
    return nil
}

type NodeValueSpec struct {
}

func (n NodeValueSpec) Type() NodeType {
    return ValueSpec
}

func (n NodeValueSpec) Children() []Node {
    return nil
}

type File struct {
    Package string
    Imports []string
    Decl []Node
}

func parsePackage(tokens []Token) (string, []Token, error) {
    if tokens[0].id != PACKAGE || tokens[1].id != IDENT {
        return "", tokens, errors.New("expected package statement first")
    }
    return tokens[1].name, tokens[2:], nil
}

func parseImports(tokens []Token) ([]string, []Token, error) {
    var imports []string
    if tokens[0].id != IMPORT {
        return imports, tokens, nil //it's okay to not import anything
    }
    if tokens[1].id == STRING {
       imports = append(imports, tokens[1].name)
       return imports, tokens[2:], nil
    }

    pos := 2
    if tokens[1].id == LPAREN {
        for tokens[pos].id != RPAREN {
            if tokens[pos].id != STRING {
                return imports, tokens, errors.New("expected string with imported package")
            }
            imports = append(imports, tokens[pos].name)
            pos++
        }
    }
    return imports, tokens[pos:], nil
}

func parseValueSpec(tokens []Token) (NodeValueSpec, []Token, error) {
    return NodeValueSpec{}, tokens, nil
}

func parseTypeSpec(tokens []Token) (NodeTypeSpec, []Token, error) {
    return NodeTypeSpec{}, tokens, nil
}

func parseFuncDecl(tokens []Token) (NodeFunc, []Token, error) {
    return NodeFunc{name: tokens[0].name}, tokens, nil
}

func parseDeclaration(tokens []Token) (Node, []Token, error) {
    switch tokens[0].id {
    case CONST, VAR:
        return parseValueSpec(tokens[1:])

	case TYPE:
        return parseTypeSpec(tokens[1:])

	case FUNC:
        return parseFuncDecl(tokens[1:])
    }
    return NodeFunc{}, tokens, errors.New("Expected declaration")
}

func parse(tokens []Token) (File, error) {
    fmt.Println(tokens)
    file := File{}
    package_name, tokens, err := parsePackage(tokens)
    if err != nil {
        return file, err
    }
    file.Package = package_name

    imports, tokens, err := parseImports(tokens)
    if err != nil {
        return file, err
    }
    file.Imports = imports

    for tokens[0].id != EOF {
        var decl Node;
        decl, tokens, err = parseDeclaration(tokens)
        if err != nil {
            return file, err
        }
        file.Decl = append(file.Decl, decl)

        // Debug
        return file, nil
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
