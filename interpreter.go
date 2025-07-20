package main

import (
	"fmt"
	"io"
	"iter"
	"os"
	"strings"

	"github.com/JoshPattman/jpf"
)

type Scope struct {
	variableLevels []map[string]string
	funcitonLevels []map[string]FuncDefNode
}

func NewScope() *Scope {
	return &Scope{
		variableLevels: []map[string]string{
			{},
		},
		funcitonLevels: []map[string]FuncDefNode{
			{},
		},
	}
}

func (s *Scope) Set(key, val string) {
	for _, level := range s.variableLevels {
		if _, ok := level[key]; ok {
			level[key] = val
			return
		}
	}
	s.variableLevels[len(s.variableLevels)-1][key] = val
}

func (s *Scope) Get(key string) string {
	for _, level := range s.variableLevels {
		if _, ok := level[key]; ok {
			return level[key]
		}
	}
	panic("not in scope")
}

func (s *Scope) Has(key string) bool {
	for _, level := range s.variableLevels {
		if _, ok := level[key]; ok {
			return true
		}
	}
	return false
}

func (s *Scope) Del(key string) {
	for _, level := range s.variableLevels {
		if _, ok := level[key]; ok {
			delete(level, key)
			return
		}
	}
	panic("not in scope")
}

// Function support

func (s *Scope) SetFunc(key string, fn FuncDefNode) {
	for _, level := range s.funcitonLevels {
		if _, ok := level[key]; ok {
			level[key] = fn
			return
		}
	}
	s.funcitonLevels[len(s.funcitonLevels)-1][key] = fn
}

func (s *Scope) GetFunc(key string) FuncDefNode {
	for _, level := range s.funcitonLevels {
		if fn, ok := level[key]; ok {
			return fn
		}
	}
	panic("function not in scope")
}

func (s *Scope) HasFunc(key string) bool {
	for _, level := range s.funcitonLevels {
		if _, ok := level[key]; ok {
			return true
		}
	}
	return false
}

func (s *Scope) DelFunc(key string) {
	for _, level := range s.funcitonLevels {
		if _, ok := level[key]; ok {
			delete(level, key)
			return
		}
	}
	panic("function not in scope")
}

func (s *Scope) SubScope() *Scope {
	newVarLevels := make([]map[string]string, len(s.variableLevels)+1)
	copy(newVarLevels, s.variableLevels)
	newVarLevels[len(newVarLevels)-1] = make(map[string]string)

	newFuncLevels := make([]map[string]FuncDefNode, len(s.funcitonLevels)+1)
	copy(newFuncLevels, s.funcitonLevels)
	newFuncLevels[len(newFuncLevels)-1] = make(map[string]FuncDefNode)

	return &Scope{
		variableLevels: newVarLevels,
		funcitonLevels: newFuncLevels,
	}
}

func (s *Scope) KVPs() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for _, l := range s.variableLevels {
			for k, v := range l {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Copies the functions to the innermost level in s
func (s *Scope) CopyFuncsFrom(other *Scope) {
	for _, level := range other.funcitonLevels {
		for k, v := range level {
			s.funcitonLevels[len(s.funcitonLevels)-1][k] = v
		}
	}
}

func BuildIntereterModel() jpf.Model {
	url := os.Getenv("OPENAI_URL")
	modelName := os.Getenv("OPENAI_MODEL")
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}
	builder := jpf.BuildOpenAIModel(
		os.Getenv("OPENAI_KEY"),
		modelName,
		false,
	)
	if url != "" {
		builder = builder.WithURL(url)
	}
	model, err := builder.Validate()
	if err != nil {
		panic(err)
	}
	return model
}

func Interpret(code []ASTNode, args []string, stdout io.Writer) error {
	_, err := interpret(code, args, stdout, NewScope())
	return err
}

// If an interpret returns a non-nil value list, a return has been triggered and needs to be caught by a function. It will propagate.
func interpret(code []ASTNode, args []string, stdout io.Writer, scope *Scope) ([]string, error) {
	for _, node := range code {
		if vals, err := interpretNode(node, args, stdout, scope); err != nil {
			return nil, err
		} else if vals != nil {
			return vals, nil
		}
	}
	return nil, nil
}

func interpretNode(code ASTNode, args []string, stdout io.Writer, scope *Scope) ([]string, error) {
	switch code := code.(type) {
	case LetNode:
		err := interpretLet(code, scope)
		return nil, err
	case UseNode:
		err := interpretUse(code, scope, args)
		return nil, err
	case IfNode:
		return interpretIf(code, scope, args, stdout)
	case WhileNode:
		return interpretWhile(code, scope, args, stdout)
	case PrintNode:
		err := interpretPrint(code, scope, stdout)
		return nil, err
	case CommentNode:
		err := interpretComment(code, scope)
		return nil, err
	case DelNode:
		err := interpretDel(code, scope)
		return nil, err
	case FuncDefNode:
		err := interpretFuncDef(code, scope)
		return nil, err
	case RunNode:
		return interpretRun(code, args, stdout, scope)
	case ReturnNode:
		return interpretReturn(code, scope)
	default:
		panic(fmt.Sprintf("unrecognised node type %T", code))
	}
}

func interpretLet(n LetNode, scope *Scope) error {
	prompt := "You have been asked to set the value of a variable in an LLM-based programming language." +
		"The user will ask you what to put in your anser" +
		"Your entire response will be copied verbatim into the variable value. For this reason, you don't need to specity code to set the variable (e.g. omit." +
		"Current other variables in scope at the moment are:\n" +
		" $SCOPE$"

	scopeVars := []string{}
	for k, v := range scope.KVPs() {
		scopeVars = append(scopeVars, "## VARIABLE "+k+"\n"+v)
	}
	scopeVarsStr := strings.Join(scopeVars, "\n\n")
	prompt = strings.ReplaceAll(prompt, "$SCOPE$", scopeVarsStr)

	model := BuildIntereterModel()
	_, resp, _, err := model.Respond([]jpf.Message{
		{Role: jpf.SystemRole, Content: prompt},
		{Role: jpf.UserRole, Content: n.Value},
	})
	if err != nil {
		return fmt.Errorf("error interpreting let node: %w", err)
	}
	scope.Set(n.Ident, strings.TrimSpace(resp.Content))
	return nil
}

func interpretUse(n UseNode, scope *Scope, args []string) error {
	if n.ArgID >= 0 && n.ArgID < len(args) {
		scope.Set(n.Ident, args[n.ArgID])
		return nil
	} else {
		return fmt.Errorf("argument id %d is out or range for arguments", n.ArgID)
	}
}

func interpretIf(n IfNode, scope *Scope, args []string, out io.Writer) ([]string, error) {
	prompt := "You have been asked to evaluate the truthyness of a statement in an LLM-based programming language." +
		"The user will ask you what to put in your anser" +
		"You can use variables in scope to give your answer context." +
		"Your response MUST eventually contain either 'EVALUATE_TRUE' or 'EVALUATE_FALSE'." +
		"Current other variables in scope at the moment are:\n" +
		" $SCOPE$"

	scopeVars := []string{}
	for k, v := range scope.KVPs() {
		scopeVars = append(scopeVars, "## VARIABLE "+k+"\n"+v)
	}
	scopeVarsStr := strings.Join(scopeVars, "\n\n")
	prompt = strings.ReplaceAll(prompt, "$SCOPE$", scopeVarsStr)

	model := BuildIntereterModel()
	_, resp, _, err := model.Respond([]jpf.Message{
		{Role: jpf.SystemRole, Content: prompt},
		{Role: jpf.UserRole, Content: n.Condition},
	})
	if err != nil {
		return nil, fmt.Errorf("error interpreting let node: %w", err)
	}
	subScope := scope.SubScope()
	if strings.Contains(resp.Content, "EVALUATE_TRUE") {
		return interpret(n.IfStatements, args, out, subScope)
	} else if strings.Contains(resp.Content, "EVALUATE_FALSE") {
		return interpret(n.ElseStatements, args, out, subScope)
	} else {
		return nil, fmt.Errorf("llm did not decide")
	}
}

func interpretPrint(n PrintNode, scope *Scope, out io.Writer) error {
	if !scope.Has(n.Ident) {
		return fmt.Errorf("variable %s not in scope", n.Ident)
	}
	_, err := fmt.Fprintln(out, scope.Get(n.Ident))
	return err
}

func interpretComment(_ CommentNode, _ *Scope) error {
	return nil
}

func interpretWhile(n WhileNode, scope *Scope, args []string, out io.Writer) ([]string, error) {
	for {
		prompt := "You have been asked to evaluate the truthyness of a statement in an LLM-based programming language." +
			"The user will ask you what to put in your anser" +
			"You can use variables in scope to give your answer context." +
			"Your response MUST eventually contain either 'EVALUATE_TRUE' or 'EVALUATE_FALSE'." +
			"Current other variables in scope at the moment are:\n" +
			" $SCOPE$"

		scopeVars := []string{}
		for k, v := range scope.KVPs() {
			scopeVars = append(scopeVars, "## VARIABLE "+k+"\n"+v)
		}
		scopeVarsStr := strings.Join(scopeVars, "\n\n")
		prompt = strings.ReplaceAll(prompt, "$SCOPE$", scopeVarsStr)

		model := BuildIntereterModel()
		_, resp, _, err := model.Respond([]jpf.Message{
			{Role: jpf.SystemRole, Content: prompt},
			{Role: jpf.UserRole, Content: n.Condition},
		})
		if err != nil {
			return nil, fmt.Errorf("error interpreting let node: %w", err)
		}
		subScope := scope.SubScope()
		if strings.Contains(resp.Content, "EVALUATE_TRUE") {
			returnVals, err := interpret(n.Statements, args, out, subScope)
			if err != nil {
				return nil, err
			}
			// If the loop body returned values, propagate them up
			if returnVals != nil {
				return returnVals, nil
			}
		} else if strings.Contains(resp.Content, "EVALUATE_FALSE") {
			return nil, nil
		} else {
			return nil, fmt.Errorf("llm did not decide")
		}
	}
}

func interpretDel(n DelNode, scope *Scope) error {
	if !scope.Has(n.Ident) {
		return fmt.Errorf("variable %s not in scope", n.Ident)
	}
	scope.Del(n.Ident)
	return nil
}

func interpretFuncDef(n FuncDefNode, scope *Scope) error {
	scope.SetFunc(n.Ident, n)
	return nil
}

func interpretReturn(n ReturnNode, scope *Scope) ([]string, error) {
	vals := make([]string, 0)
	for _, ident := range n.Idents {
		if !scope.Has(ident) {
			return nil, fmt.Errorf("scope does not contain variable %s", ident)
		}
		vals = append(vals, scope.Get(ident))
	}
	return vals, nil
}

func interpretRun(n RunNode, args []string, stdout io.Writer, scope *Scope) ([]string, error) {
	if !scope.HasFunc(n.FnIdent) {
		return nil, fmt.Errorf("function %s is not defined", n.FnIdent)
	}
	fn := scope.GetFunc(n.FnIdent)
	if len(fn.Args) != len(n.InputIdents) {
		return nil, fmt.Errorf("function expected %d args but got %d", len(fn.Args), len(n.InputIdents))
	}
	freshScope := NewScope()
	for i, ident := range n.InputIdents {
		if !scope.Has(ident) {
			return nil, fmt.Errorf("variable %s not in scope", ident)
		}
		freshScope.Set(fn.Args[i], scope.Get(ident))
	}
	freshScope.CopyFuncsFrom(scope)
	returnVal, err := interpret(fn.Code, args, stdout, freshScope)
	if err != nil {
		return nil, err
	}
	if len(returnVal) != len(n.OutputIdents) {
		return nil, fmt.Errorf("function provided %d outputs but caller provides %d", len(returnVal), len(n.OutputIdents))
	}
	for i, ident := range n.OutputIdents {
		scope.Set(ident, returnVal[i])
	}
	return nil, nil
}
