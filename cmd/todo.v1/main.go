package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gregidonut/basic_todo_app"
	"io"
	"os"
	"strings"
)

// Default file name
var todoFileName = ".todo.json"

func main() {
	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Parsing command line flags

	// these flags doesn't need to return val since we will use container to
	// usedFlag map to check for presence which is already a bool
	flag.Bool("add", false, "Task to be included in the ToDo list")
	flag.Bool("list", false, "List all tasks")

	complete := flag.Int("complete", 0, "Item to be completed")
	deletedItem := flag.Int("del", 0, "Item to be deleted")

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
	case usedFlag["add"]:
		// When any arguments (excluding flags) are provided, they will be
		// used as the new task
		task, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.Add(task)

		err = l.Save(todoFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case usedFlag["del"]:
		// deleted the given item
		derefList := *l
		itemToBeDeleted := derefList[*deletedItem-1].Task

		err := l.Delete(*deletedItem)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Deleted task number %d: %s\n", *deletedItem, itemToBeDeleted)

		// Save the new list
		err = l.Save(todoFileName)
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

// getTask function decides where to get the description for a new task
// from: arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()

	err := s.Err()
	if err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank\n")
	}

	return s.Text(), nil
}
