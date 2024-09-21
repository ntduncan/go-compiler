package main

import (
	"fmt"
	"ntduncan.com/go-compiler/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Get started by typing a command\n")
	repl.Start(os.Stdin, os.Stdout)
}
