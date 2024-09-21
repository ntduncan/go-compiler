package repl

import (
	"bufio"
	"fmt"
	"io"

	"ntduncan.com/go-compiler/evaluator"
	"ntduncan.com/go-compiler/lexer"
	"ntduncan.com/go-compiler/object"
	"ntduncan.com/go-compiler/parser"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")

		}

	}
}

const MONKEY_FACE = `
	 --,--
 .--.  .-"   "-.  .--.
/ .. \/ .-. .-. \/ .. \
| | '| /   Y   \ |'  | |
| \  \ \ 0 | 0 / /   / |
\ '-,\.-"""""""-./,-' /
 ''-'/_   ^ ^   _\'-''
     |  \._ _./  |
      \ \ '~' / /
      '._'-=-'_.'
	'-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some munkey business here!\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}