package main

import (
    "fmt"
    "io/ioutil"
    "errors"
    "regexp"
)

func match(pattern string, corpse string, pos int)  []string {
    re := regexp.MustCompile(pattern)
    return re.FindStringSubmatch(corpse[pos:])
}

type TokenID int

type Token struct {
    id TokenID
    name string
}

const (
    ILLEGAL TokenID = iota
    EOF
	COMMENT

    IDENT  // main
	INT    // 12345
    HEX    // 0x2
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	RAWSTRING // `abcas`

    ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

    keyword_beg
	// Keywords
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR

	FUNC
	GO
	GOTO
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	BREAK:    "break",
	CASE:     "case",
	CHAN:     "chan",
	CONST:    "const",
	CONTINUE: "continue",

	DEFAULT:     "default",
	DEFER:       "defer",
	ELSE:        "else",
	FALLTHROUGH: "fallthrough",
	FOR:         "for",

	FUNC:   "func",
	GO:     "go",
	GOTO:   "goto",
	IF:     "if",
	IMPORT: "import",

	INTERFACE: "interface",
	MAP:       "map",
	PACKAGE:   "package",
	RANGE:     "range",
	RETURN:    "return",

	SELECT: "select",
	STRUCT: "struct",
	SWITCH: "switch",
	TYPE:   "type",
	VAR:    "var",
}

func (tok Token) String() string {
    if len(tok.name) > 0 {
        return tokens[tok.id] + "(" + tok.name + ")"
    }
    return tokens[tok.id]
}


var keywords map[string]TokenID

func MakeKeywordMap() {
	keywords = make(map[string]TokenID)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) TokenID {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}


func scan(code string, pos *int) (Token, error) {
    // skip withspace
    for code[*pos] == ' ' || code[*pos] == '\t' || code[*pos] == '\n' {
        *pos += 1
        if *pos >= len(code) {
            return Token{EOF,""}, nil
        }
    }
    // scan comments
    if m := match(`^\/\/[^\n]*\n`, code, *pos); len(m) > 0 {
        // Matches single line comments
       *pos += len(m[0])
       return Token{COMMENT,""}, nil
    } else if code[*pos] == '(' {
       *pos++
       return Token{LPAREN,""}, nil
    } else if code[*pos] == ')' {
       *pos++
       return Token{RPAREN,""}, nil
    } else if code[*pos] == '{' {
       *pos++
       return Token{LBRACE,""}, nil
    } else if code[*pos] == '}' {
       *pos++
       return Token{RBRACE,""}, nil
    } else if code[*pos] == '.' {
       *pos++
       return Token{PERIOD,""}, nil
    } else if m := match(`^\/\*[^\*\/]*(\*\/)?`, code, *pos) ; len(m) > 0 {
       // Matches multi-line comments and detects if not closed.
       if len(m[1]) == 0 {
            return Token{ILLEGAL,""}, errors.New(fmt.Sprintf("Multi-line comment never closes."))
        }
       *pos += len(m[0])
       return Token{COMMENT,""}, nil
    } else if m := match(`^([_a-zA-Z][_a-zA-Z0-9]*)`, code, *pos); len(m) > 0 {
      // find identifier
       *pos += len(m[0])
       return Token{Lookup(m[1]), m[1]}, nil
    } else if m := match(`^(\".*\")`, code, *pos); len(m) > 0 {
       // scan string
       *pos += len(m[0])
       return Token{STRING, m[1]}, nil
    } else if m := match(`^\x60([^\x60]*)\x60`, code, *pos); len(m) > 0 {
       // scan string
       *pos += len(m[0])
       return Token{RAWSTRING, m[1]}, nil
    } else if m := match(`^((0x)?(\d*)(\.)?(\d*)(i)?)`, code, *pos); len(m) > 0 {
       // scan number
       *pos += len(m[0])
       return Token{INT, m[1]}, nil
    } else if m := match(`^\.\.\.`, code, *pos); len(m) > 0 {
       // scan number
       *pos += len(m[0])
       return Token{ELLIPSIS,""}, nil
    } else if code[*pos] == ':' {
       *pos++
       return Token{COLON, ""}, nil
    } else {
        panic(fmt.Sprintf("Unhandled %s", code[*pos:*pos+5]))
        return Token{ILLEGAL,""}, errors.New(fmt.Sprintf("Unhandled %s", code[*pos:*pos+5]))
    }

    // elipsis
    // ,
    // ;
    // (
    // )
    // [
    // ]
    // {
    // }
    // +
    // -
    // *
    // /
    // %
    // ^
    // 
    return Token{EOF,""}, nil
}

func tokenize(code string) ([]Token, error) {
    MakeKeywordMap()
    var tokens []Token
    pos := 0
    var token Token
    var err error
    for token.id != EOF {
        token, err = scan(code, &pos)
        tokens = append(tokens, token)
        if err != nil {
            panic(err)
        }
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

