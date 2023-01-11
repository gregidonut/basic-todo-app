package main

import (
	"fmt"
	"github.com/gregidonut/basic_todo_app"
	"os"
)

// temporarily hard-coding file name here
const todoFileName = ".todo.json"

func main() {
	l := &basic_todo_app.List{}

	// Use the Get method to read to do items from file
	err := l.Get(todoFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
