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
	// this flag doesn't need to return val since we will use container to
	// usedFlag map to check for presence which is already a bool
	flag.Bool("list", false, "List all tasks")
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
	case usedFlag["list"]:
		listToDoItems(l)
	case usedFlag["complete"]:
		// Complete the given item
		err := l.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		err = l.Save(todoFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case usedFlag["task"]:
		l.Add(*task)

		err := l.Save(todoFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case len(os.Args) == 1:
		// no arguments provided
		listToDoItems(l)
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// extracting this logic to make switch statement more readable
// by minimizing indents
func listToDoItems(l *basic_todo_app.List) {
	// List current to do items if itemsList is not empty
	if len(*l) <= 0 {
		fmt.Println("You have no to do items")
		os.Exit(0)
	}
	fmt.Print(l)
}
