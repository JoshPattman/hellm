package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ASTNode interface {
	Format(indent string) string
}

type LetNode struct {
	Ident string
	Value string
}

type UseNode struct {
	Ident string
	ArgID int
}

type IfNode struct {
	Condition      string
	IfStatements   []ASTNode
	ElseStatements []ASTNode
}

type WhileNode struct {
	Condition  string
	Statements []ASTNode
}

type PrintNode struct {
	Ident string
}

type CommentNode struct {
	Comment string
}

type DelNode struct {
	Ident string
}

type RunNode struct {
	OutputIdents []string
	FnIdent      string
	InputIdents  []string
}

type ReturnNode struct {
	Idents []string
}

type FuncDefNode struct {
	Ident string
	Args  []string
	Code  []ASTNode
}

func (n LetNode) Format(indent string) string {
	return fmt.Sprintf("%slet %s = \"%s\";", indent, n.Ident, n.Value)
}
func (n UseNode) Format(indent string) string {
	return fmt.Sprintf("%suse %s = %d;", indent, n.Ident, n.ArgID)
}
func (n IfNode) Format(indent string) string {
	stmtFormats := make([]string, len(n.IfStatements))
	for i, stmt := range n.IfStatements {
		stmtFormats[i] = stmt.Format(indent + "    ")
	}
	elseStmtFormats := make([]string, len(n.ElseStatements))
	for i, stmt := range n.ElseStatements {
		elseStmtFormats[i] = stmt.Format(indent + "    ")
	}
	if len(elseStmtFormats) == 0 {
		return fmt.Sprintf("%sif \"%s\" {\n%v\n%s}", indent, n.Condition, strings.Join(stmtFormats, "\n"), indent)
	} else {
		return fmt.Sprintf("%sif \"%s\" {\n%v\n%s} else {\n%s\n%s}", indent, n.Condition, strings.Join(stmtFormats, "\n"), indent, strings.Join(elseStmtFormats, "\n"), indent)
	}
}
func (n WhileNode) Format(indent string) string {
	stmtFormats := make([]string, len(n.Statements))
	for i, stmt := range n.Statements {
		stmtFormats[i] = stmt.Format(indent + "    ")
	}
	return fmt.Sprintf("%swhile \"%s\" {\n%s\n%s}", indent, n.Condition, strings.Join(stmtFormats, "\n"), indent)
}
func (n PrintNode) Format(indent string) string {
	return fmt.Sprintf("%sprint %s;", indent, n.Ident)
}
func (n CommentNode) Format(indent string) string {
	return fmt.Sprintf("\n%scom \"%s\";", indent, n.Comment)
}
func (n DelNode) Format(indent string) string {
	return fmt.Sprintf("%sdel %s;", indent, n.Ident)
}
func (n RunNode) Format(indent string) string {
	outputs := ""
	if len(n.OutputIdents) > 0 {
		outputs = strings.Join(n.OutputIdents, " ")
	}
	inputs := ""
	if len(n.InputIdents) > 0 {
		inputs = strings.Join(n.InputIdents, " ")
	}
	if outputs != "" && inputs != "" {
		return fmt.Sprintf("%srun %s = %s %s;", indent, outputs, n.FnIdent, inputs)
	} else if outputs != "" {
		return fmt.Sprintf("%srun %s = %s;", indent, outputs, n.FnIdent)
	} else if inputs != "" {
		return fmt.Sprintf("%srun %s %s;", indent, n.FnIdent, inputs)
	} else {
		return fmt.Sprintf("%srun %s;", indent, n.FnIdent)
	}
}
func (n ReturnNode) Format(indent string) string {
	idents := strings.Join(n.Idents, " ")
	return fmt.Sprintf("%sreturn %s;", indent, idents)
}
func (n FuncDefNode) Format(indent string) string {
	args := ""
	if len(n.Args) > 0 {
		args = " " + strings.Join(n.Args, " ")
	}
	codeLines := make([]string, len(n.Code))
	for i, stmt := range n.Code {
		codeLines[i] = stmt.Format(indent + "    ")
	}
	return fmt.Sprintf("\n%sfn %s%s {\n%s\n%s}\n", indent, n.Ident, args, strings.Join(codeLines, "\n"), indent)
}

func Parse(tokens []LexToken) ([]ASTNode, error) {
	nodes := []ASTNode{}
	for len(tokens) > 0 {
		node, rest, ok := tryParseNode(tokens)
		if !ok {
			return nil, fmt.Errorf("failed to parse tokens: %v", FormatLexTokens(tokens))
		}
		nodes = append(nodes, node)
		tokens = rest
	}

	if len(tokens) > 0 {
		return nil, fmt.Errorf("unexpected tokens remaining after parsing: %v", tokens)
	}

	return nodes, nil

}

type PatternMatchable interface {
	Copy(tokens []LexToken) (int, bool)
}

type patternMatchList[T LexToken] struct {
	elems []T
}

// Copies tokens from the tokens array into our array until a token that is the wrong type is reached.
// Returns the amount it copied.
func (p *patternMatchList[T]) Copy(tokens []LexToken) (int, bool) {
	for i, tok := range tokens {
		if tok, ok := tok.(T); !ok {
			return i, true
		} else {
			p.elems = append(p.elems, tok)
		}
	}
	return len(tokens), true
}

func patternMatch(tokens []LexToken, pattern ...PatternMatchable) (bool, []LexToken) {
	ti := 0
	for _, pat := range pattern {
		if ti > len(tokens) {
			return false, nil
		}
		inc, ok := pat.Copy(tokens[ti:])
		if !ok {
			return false, nil
		}
		ti += inc
	}
	return true, tokens[ti:]
}

func parseNodesUntilNoMoreParse(tokens []LexToken) ([]ASTNode, []LexToken) {
	nodes := []ASTNode{}
	for len(tokens) > 0 {
		node, rest, ok := tryParseNode(tokens)
		if !ok {
			break
		}
		nodes = append(nodes, node)
		tokens = rest
	}
	return nodes, tokens
}

func tryParseNode(tokens []LexToken) (ASTNode, []LexToken, bool) {
	parseFuncs := []func([]LexToken) (ASTNode, []LexToken, bool){
		tryParseLet,
		tryParseUse,
		tryParseIf,
		tryParseFuncDef,
		tryParseWhile,
		tryParsePrint,
		tryParseDel,
		tryParseComment,
		tryParseReturn,
		tryParseRun,
		tryParseNoAssnRun,
	}
	for _, parseFunc := range parseFuncs {
		node, rest, ok := parseFunc(tokens)
		if ok {
			return node, rest, true
		}
	}
	return nil, tokens, false
}

func tryParseDel(tokens []LexToken) (ASTNode, []LexToken, bool) {
	ident := &IdentLexToken{}
	if ok, rest := patternMatch(tokens, &DelLexToken{}, ident, &SemiColonLexToken{}); ok {
		return DelNode{
			Ident: ident.Name,
		}, rest, true
	}
	return nil, nil, false
}

func tryParseLet(tokens []LexToken) (ASTNode, []LexToken, bool) {
	ident := &IdentLexToken{}
	value := &StringLexToken{}
	if ok, rest := patternMatch(tokens, &LetLexToken{}, ident, &EqLexToken{}, value, &SemiColonLexToken{}); ok {
		return LetNode{
			Ident: ident.Name,
			Value: value.Value,
		}, rest, true
	}
	return nil, nil, false
}

func tryParseUse(tokens []LexToken) (ASTNode, []LexToken, bool) {
	ident := &IdentLexToken{}
	argID := &IdentLexToken{}
	if ok, rest := patternMatch(tokens, &UseLexToken{}, ident, &EqLexToken{}, argID, &SemiColonLexToken{}); ok {
		id, err := strconv.Atoi(argID.Name)
		if err != nil {
			return nil, nil, false
		}
		return UseNode{
			Ident: ident.Name,
			ArgID: id,
		}, rest, true
	}
	return nil, nil, false
}

func tryParseIf(tokens []LexToken) (ASTNode, []LexToken, bool) {
	condition := &StringLexToken{}
	if ok, tokens := patternMatch(tokens, &IfLexToken{}, condition, &OpenBraceLexToken{}); ok {
		var ifChildren, elseChildren []ASTNode
		ifChildren, tokens = parseNodesUntilNoMoreParse(tokens)
		if ok, tokens = patternMatch(tokens, &CloseBraceLexToken{}); !ok {
			return nil, nil, false
		}
		if ok, elseTokens := patternMatch(tokens, &ElseLexToken{}, &OpenBraceLexToken{}); ok {
			elseChildren, elseTokens = parseNodesUntilNoMoreParse(elseTokens)
			if ok, elseTokens = patternMatch(elseTokens, &CloseBraceLexToken{}); !ok {
				return nil, nil, false
			}
			tokens = elseTokens
		}
		return IfNode{
			Condition:      condition.Value,
			IfStatements:   ifChildren,
			ElseStatements: elseChildren,
		}, tokens, true
	}
	return nil, nil, false
}

func tryParseFuncDef(tokens []LexToken) (ASTNode, []LexToken, bool) {
	name := &IdentLexToken{}
	args := &patternMatchList[*IdentLexToken]{}
	if ok, tokens := patternMatch(tokens, &FnLexToken{}, name, args, &OpenBraceLexToken{}); ok {
		var children []ASTNode
		children, tokens = parseNodesUntilNoMoreParse(tokens)
		if ok, tokens = patternMatch(tokens, &CloseBraceLexToken{}); !ok {
			return nil, nil, false
		}
		argNames := make([]string, len(args.elems))
		for i, arg := range args.elems {
			argNames[i] = arg.Name
		}
		return FuncDefNode{
			Ident: name.Name,
			Args:  argNames,
			Code:  children,
		}, tokens, true
	}
	return nil, nil, false
}

func tryParseWhile(tokens []LexToken) (ASTNode, []LexToken, bool) {
	condition := &StringLexToken{}
	if ok, tokens := patternMatch(tokens, &WhileLexToken{}, condition, &OpenBraceLexToken{}); ok {
		statements, tokens := parseNodesUntilNoMoreParse(tokens)
		if ok, tokens := patternMatch(tokens, &CloseBraceLexToken{}); ok {
			return WhileNode{
				Condition:  condition.Value,
				Statements: statements,
			}, tokens, true
		}
	}
	return nil, nil, false
}

func tryParsePrint(tokens []LexToken) (ASTNode, []LexToken, bool) {
	message := &IdentLexToken{}
	if ok, rest := patternMatch(tokens, &PrintLexToken{}, message, &SemiColonLexToken{}); ok {
		return PrintNode{
			Ident: message.Name,
		}, rest, true
	}
	return nil, nil, false
}

func tryParseComment(tokens []LexToken) (ASTNode, []LexToken, bool) {
	message := &StringLexToken{}
	if ok, rest := patternMatch(tokens, &CommentLexToken{}, message, &SemiColonLexToken{}); ok {
		return CommentNode{
			Comment: message.Value,
		}, rest, true
	}
	return nil, nil, false
}

func tryParseReturn(tokens []LexToken) (ASTNode, []LexToken, bool) {
	idents := &patternMatchList[*IdentLexToken]{}
	if ok, rest := patternMatch(tokens, &ReturnLexToken{}, idents, &SemiColonLexToken{}); ok {
		identNames := make([]string, len(idents.elems))
		for i, ident := range idents.elems {
			identNames[i] = ident.Name
		}
		return ReturnNode{
			Idents: identNames,
		}, rest, true
	}
	return nil, nil, false
}

func tryParseRun(tokens []LexToken) (ASTNode, []LexToken, bool) {
	outputIdents := &patternMatchList[*IdentLexToken]{}
	inputIdents := &patternMatchList[*IdentLexToken]{}
	fnIdent := &IdentLexToken{}
	if ok, rest := patternMatch(tokens, &RunLexToken{}, outputIdents, &EqLexToken{}, fnIdent, inputIdents, &SemiColonLexToken{}); ok {
		outputs := make([]string, len(outputIdents.elems))
		for i, ident := range outputIdents.elems {
			outputs[i] = ident.Name
		}
		inputs := make([]string, len(inputIdents.elems))
		for i, ident := range inputIdents.elems {
			inputs[i] = ident.Name
		}
		return RunNode{
			OutputIdents: outputs,
			InputIdents:  inputs,
			FnIdent:      fnIdent.Name,
		}, rest, true
	}
	return nil, nil, false
}

func tryParseNoAssnRun(tokens []LexToken) (ASTNode, []LexToken, bool) {
	inputIdents := &patternMatchList[*IdentLexToken]{}
	fnIdent := &IdentLexToken{}
	if ok, rest := patternMatch(tokens, &RunLexToken{}, fnIdent, inputIdents, &SemiColonLexToken{}); ok {
		inputs := make([]string, len(inputIdents.elems))
		for i, ident := range inputIdents.elems {
			inputs[i] = ident.Name
		}
		return RunNode{
			OutputIdents: []string{},
			InputIdents:  inputs,
			FnIdent:      fnIdent.Name,
		}, rest, true
	}
	return nil, nil, false
}
