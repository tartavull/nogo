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

func parseArguments(tokens []Token) ([]string, []string, []Token, error) {
    if (tokens[0].id != LPAREN) {
        return nil, nil, tokens, errors.New("Expected (")
    }
    if (tokens[1].id != RPAREN) {
        return nil, nil, tokens, errors.New("Expected )")
    }
    return []string{}, []string{}, tokens[2:], nil
}

func parseSignature(tokens []Token) ([]string, []Token, error) {
    return []string{}, tokens, nil
}

func parseStmt(tokens []Token) {
    switch tokens[0].id {
    case CONST, TYPE, VAR:
        panic("not implemented")
    case
		// tokens that may start an expression
		IDENT, INT, FLOAT, IMAG, CHAR, STRING, FUNC, LPAREN, // operands
		LBRACK, STRUCT, MAP, CHAN, INTERFACE, // composite types
		ADD, SUB, MUL, AND, XOR, ARROW, NOT: // unary operators
		//s, _ = p.parseSimpleStmt(labelOk)
		// because of the required look-ahead, labeled statements are
		// parsed by parseSimpleStmt - don't expect a semicolon after
		// them
        panic("not implemented")
	case GO:
		//s = p.parseGoStmt()
        panic("not implemented")
	case DEFER:
		//s = p.parseDeferStmt()
        panic("not implemented")
	case RETURN:
		//s = p.parseReturnStmt()
        panic("not implemented")
	case BREAK, CONTINUE, GOTO, FALLTHROUGH:
		//s = p.parseBranchStmt(p.tok)
	case LBRACE:
        panic("not implemented")
		//s = p.parseBlockStmt()
		//p.expectSemi()
	case IF:
        panic("not implemented")
		//s = p.parseIfStmt()
	case SWITCH:
        panic("not implemented")
		//s = p.parseSwitchStmt()
	case SELECT:
        panic("not implemented")
		//s = p.parseSelectStmt()
	case FOR:
        panic("not implemented")
		//s = p.parseForStmt()
	case SEMICOLON:
        panic("not implemented")
		// Is it ever possible to have an implicit semicolon
		// producing an empty statement in a valid program?
		// (handle correctly anyway)
		//s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: p.lit == "\n"}
		//p.next()
	case RBRACE:
        panic("not implemented")
		// a semicolon may be omitted before a closing "}"
		//s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: true}
    }
}

func parseBody(tokens []Token) ([]Node, []Token, error) {
    if (tokens[0].id != LBRACE) {
        return nil, tokens, errors.New("Expected {")
    }
    parseStmt(tokens[1:])
    return []Node{}, tokens, nil

}

func parseFuncDecl(tokens []Token) (NodeFunc, []Token, error) {
    var err error
    decl := NodeFunc{}
    if (tokens[0].id != IDENT) {
        return NodeFunc{}, tokens, errors.New("Expected function identifier")
    }
    decl.name = tokens[0].name
    decl.arg_name, decl.arg_type, tokens, err = parseArguments(tokens[1:])
    if (err != nil) {
        return NodeFunc{}, tokens, err
    }

    decl.ret_type, tokens, err = parseSignature(tokens)
    if (err != nil) {
        return NodeFunc{}, tokens, err
    }
    decl.body, tokens, err = parseBody(tokens)
    if (err != nil) {
        return NodeFunc{}, tokens, err
    }

    return decl, tokens, nil
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
