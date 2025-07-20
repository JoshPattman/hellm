package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		failf("please provide a subcommand")
	}
	subCommand, commandArgs := os.Args[1], os.Args[2:]

	switch subCommand {
	case "help":
		printUsage()
	case "run":
		err := cmdRun(commandArgs)
		if err != nil {
			fail(err)
		}
	case "parse":
		err := cmdParse(commandArgs)
		if err != nil {
			fail(err)
		}
	case "tokenize":
		err := cmdTokenize(commandArgs)
		if err != nil {
			fail(err)
		}
	case "format":
		err := cmdFormat(commandArgs)
		if err != nil {
			fail(err)
		}
	default:
		printUsage()
		failf("unrecognised subcommand '%s'", subCommand)
	}
}

func cmdRun(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("must provide a filename to run")
	}
	fileName := args[0]
	content, err := readFile(fileName)
	if err != nil {
		fail(err)
	}
	lexTokens, err := Lex(content)
	if err != nil {
		fail(err)
	}

	parsed, err := Parse(lexTokens)
	if err != nil {
		fail(err)
	}

	err = Interpret(parsed, args[1:], os.Stdout)
	if err != nil {
		fail(err)
	}

	return nil
}

func cmdParse(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("must provide a filename to run")
	}
	fileName := args[0]
	content, err := readFile(fileName)
	if err != nil {
		fail(err)
	}
	lexTokens, err := Lex(content)
	if err != nil {
		fail(err)
	}

	parsed, err := Parse(lexTokens)
	if err != nil {
		fail(err)
	}

	for _, node := range parsed {
		fmt.Println(node.Format(""))
	}

	return nil
}

func cmdTokenize(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("must provide a filename to run")
	}
	fileName := args[0]
	content, err := readFile(fileName)
	if err != nil {
		fail(err)
	}
	lexTokens, err := Lex(content)
	if err != nil {
		fail(err)
	}

	fmt.Println(FormatLexTokens(lexTokens))

	return nil
}

func cmdFormat(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("must provide a filename to run")
	}
	fileName := args[0]
	content, err := readFile(fileName)
	if err != nil {
		fail(err)
	}
	lexTokens, err := Lex(content)
	if err != nil {
		fail(err)
	}

	parsed, err := Parse(lexTokens)
	if err != nil {
		fail(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, node := range parsed {
		fmt.Fprintln(f, node.Format(""))
	}

	return nil
}

func readFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("source file '%s' does not exist", fileName)
	} else if err != nil {
		return "", errors.Join(fmt.Errorf("error reading source file '%s'", fileName), err)
	}
	return string(data), nil
}

func printUsage() {
	fmt.Println("hellm - A language for 100x devs")
	fmt.Println("usage:")
	fmt.Println("$ hellm run <filename>")
	fmt.Println("$ hellm tokenize <filename>")
	fmt.Println("$ hellm parse <filename>")
	fmt.Println("$ hellm format <filename>")
	fmt.Println("$ hellm help")

}

func fail(err error) {
	if err != nil {
		fmt.Println("fatal error:", err)
		os.Exit(1)
	}
}

func failf(f string, args ...any) {
	fail(fmt.Errorf(f, args...))
}
