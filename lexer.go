package main

import (
	"errors"
	"fmt"
	"strings"
)

type LexToken interface {
	PatternMatchable
}

type LetLexToken struct{}
type UseLexToken struct{}
type FnLexToken struct{}
type IdentLexToken struct {
	Name string
}
type StringLexToken struct {
	Value string
}
type OpenBraceLexToken struct{}
type CloseBraceLexToken struct{}
type SemiColonLexToken struct{}
type EqLexToken struct{}
type IfLexToken struct{}
type WhileLexToken struct{}
type ElseLexToken struct{}
type PrintLexToken struct{}
type CommentLexToken struct{}
type DelLexToken struct{}
type RunLexToken struct{}
type ReturnLexToken struct{}

func (t *LetLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*LetLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *UseLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*UseLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *FnLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*FnLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *IdentLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	otherT, ok := tokens[0].(*IdentLexToken)
	if !ok {
		return 0, false
	}
	t.Name = otherT.Name
	return 1, true
}
func (t *StringLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	otherT, ok := tokens[0].(*StringLexToken)
	if !ok {
		return 0, false
	}
	t.Value = otherT.Value
	return 1, true
}
func (t *OpenBraceLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*OpenBraceLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *CloseBraceLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*CloseBraceLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *IfLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*IfLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *PrintLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*PrintLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *EqLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*EqLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *SemiColonLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*SemiColonLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *ElseLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*ElseLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *WhileLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*WhileLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *CommentLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*CommentLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *DelLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*DelLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *RunLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*RunLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}
func (t *ReturnLexToken) Copy(tokens []LexToken) (int, bool) {
	if len(tokens) < 1 {
		return 0, false
	}
	_, ok := tokens[0].(*ReturnLexToken)
	if !ok {
		return 0, false
	}
	return 1, true
}

func Lex(input string) ([]LexToken, error) {
	tokens := []LexToken{}
	for len(input) > 0 {
		input = readToNextChar(input)
		if len(input) == 0 {
			break
		}
		token, rest, err := readLexToken(input)
		if err != nil {
			return nil, errors.Join(errors.New("error lexing input"), err)
		}
		tokens = append(tokens, token)
		input = rest
	}
	if len(input) > 0 {
		return nil, fmt.Errorf("unexpected input remaining after lexing: '%s'", input)
	}
	return tokens, nil
}

func FormatLexToken(token LexToken) string {
	const (
		purple = "\033[35m"
		yellow = "\033[33m"
		green  = "\033[32m"
		reset  = "\033[0m"
	)

	switch t := token.(type) {
	case *LetLexToken:
		return purple + "let" + reset
	case *UseLexToken:
		return purple + "use" + reset
	case *FnLexToken:
		return purple + "fn" + reset
	case *IdentLexToken:
		return yellow + t.Name + reset
	case *StringLexToken:
		return green + fmt.Sprintf("\"%s\"", t.Value) + reset
	case *OpenBraceLexToken:
		return "{"
	case *CloseBraceLexToken:
		return "}"
	case *SemiColonLexToken:
		return ";"
	case *EqLexToken:
		return "="
	case *IfLexToken:
		return purple + "if" + reset
	case *ElseLexToken:
		return purple + "else" + reset
	case *WhileLexToken:
		return purple + "while" + reset
	case *PrintLexToken:
		return purple + "print" + reset
	case *CommentLexToken:
		return purple + "com" + reset
	case *DelLexToken:
		return purple + "del" + reset
	case *RunLexToken:
		return purple + "run" + reset
	case *ReturnLexToken:
		return purple + "return" + reset
	default:
		panic(fmt.Sprintf("unknown token type: %T", t))
	}
}

func FormatLexTokens(tokens []LexToken) string {
	formatted := make([]string, len(tokens))
	for i, token := range tokens {
		formatted[i] = FormatLexToken(token)
	}
	return strings.Join(formatted, " ")
}

func readToNextChar(s string) string {
	return strings.TrimLeft(s, " \t\n\r")
}

func readLexToken(s string) (LexToken, string, error) {
	readFuncs := []func(string) (LexToken, string, bool){
		readLet,
		readUse,
		readFn,
		readPrint,
		readIf,
		readWhile,
		readElse,
		readComment,
		readDel,
		readRun,
		readReturn,
		readIdent,
		readEq,
		readString,
		readSemiColon,
		readOpenBrace,
		readCloseBrace,
	}

	for _, readFunc := range readFuncs {
		if token, rest, ok := readFunc(s); ok {
			return token, rest, nil
		}
	}
	sStart := s
	if len(sStart) > 20 {
		sStart = sStart[:20] + "..."
	}
	return nil, s, fmt.Errorf("unrecognized token at start of '%s'", sStart)
}

func readLet(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "let ") {
		s = strings.TrimPrefix(s, "let")
		return &LetLexToken{}, s, true
	}
	return nil, s, false
}

func readUse(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "use ") {
		s = strings.TrimPrefix(s, "use")
		return &UseLexToken{}, s, true
	}
	return nil, s, false
}

func readFn(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "fn ") {
		s = strings.TrimPrefix(s, "fn")
		return &FnLexToken{}, s, true
	}
	return nil, s, false
}

func readDel(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "del ") {
		s = strings.TrimPrefix(s, "del")
		return &DelLexToken{}, s, true
	}
	return nil, s, false
}

func readIdent(s string) (LexToken, string, bool) {
	buf := ""
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			break
		}
		buf += string(c)
	}
	if buf == "" {
		return nil, s, false
	}
	s = s[len(buf):]
	return &IdentLexToken{Name: buf}, s, true
}

func readRun(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "run ") {
		s = strings.TrimPrefix(s, "run")
		return &RunLexToken{}, s, true
	}
	return nil, s, false
}

func readEq(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "=") {
		s = strings.TrimPrefix(s, "=")
		return &EqLexToken{}, s, true
	}
	return nil, s, false
}

func readString(s string) (LexToken, string, bool) {
	if !strings.HasPrefix(s, "\"") {
		return nil, s, false
	}
	s = strings.TrimPrefix(s, "\"")
	buf := ""
	complete := false
	for _, c := range s {
		if c == '"' {
			complete = true
			break
		} else {
			buf += string(c)
		}
	}
	if !complete {
		return nil, s, false
	}
	s = s[len(buf)+1:] // +1 for the closing quote
	return &StringLexToken{Value: buf}, s, true
}

func readSemiColon(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, ";") {
		s = strings.TrimPrefix(s, ";")
		return &SemiColonLexToken{}, s, true
	}
	return nil, s, false
}

func readOpenBrace(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "{") {
		s = strings.TrimPrefix(s, "{")
		return &OpenBraceLexToken{}, s, true
	}
	return nil, s, false
}

func readCloseBrace(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "}") {
		s = strings.TrimPrefix(s, "}")
		return &CloseBraceLexToken{}, s, true
	}
	return nil, s, false
}

func readIf(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "if ") {
		s = strings.TrimPrefix(s, "if")
		return &IfLexToken{}, s, true
	}
	return nil, s, false
}

func readComment(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "com ") {
		s = strings.TrimPrefix(s, "com")
		return &CommentLexToken{}, s, true
	}
	return nil, s, false
}

func readElse(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "else ") {
		s = strings.TrimPrefix(s, "else")
		return &ElseLexToken{}, s, true
	}
	return nil, s, false
}

func readWhile(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "while ") {
		s = strings.TrimPrefix(s, "while")
		return &WhileLexToken{}, s, true
	}
	return nil, s, false
}

func readPrint(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "print ") {
		s = strings.TrimPrefix(s, "print")
		return &PrintLexToken{}, s, true
	}
	return nil, s, false
}

func readReturn(s string) (LexToken, string, bool) {
	if strings.HasPrefix(s, "return ") || strings.HasPrefix(s, "return;") {
		s = strings.TrimPrefix(s, "return")
		return &ReturnLexToken{}, s, true
	}
	return nil, s, false
}
