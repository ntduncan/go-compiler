package repl

import (
	"bufio"
	"fmt"
	"io"

	"ntduncan.com/go-compiler/compiler"
	"ntduncan.com/go-compiler/lexer"
	"ntduncan.com/go-compiler/object"
	"ntduncan.com/go-compiler/parser"
	"ntduncan.com/go-compiler/vm"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	//env := object.NewEnvironment()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GloblasSize)
	symbolTable := compiler.NewSymbolTable()

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

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Compilation failed: \n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGloblasStore(code, globals)

		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed:\n %s\n", err)
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
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
