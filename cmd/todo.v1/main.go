package main

import (
	"flag"
	"fmt"
	"github.com/gregidonut/basic_todo_app"
	"os"
)

// temporarily hard-coding file name here
const todoFileName = ".todo.json"

func main() {
	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Parse()

	// setup container to check if a specific flag was provided
	usedFlag := make(map[string]bool)
	// add flag if used to usedFlag map (container
	flag.Visit(func(f *flag.Flag) {
		usedFlag[f.Name] = true
	})

	// Define items List
	l := &basic_todo_app.List{}

	// Use the Get method to read to do items from file
	err := l.Get(todoFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the flag used
	switch {
	case usedFlag["task"]:
	case usedFlag["completed"]:
	case usedFlag["list"]:
	default:
	}
}
