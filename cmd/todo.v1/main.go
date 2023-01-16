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
	flag.Bool("verbose", false, "shows completed or created date")
	flag.Bool("no-complete", false, "prevents showing of completed items")

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
		listToDoItems(l, usedFlag["verbose"], usedFlag["no-complete"])
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
		tasks, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, task := range tasks {
			if task == "" {
				fmt.Fprintln(os.Stderr, "task cannot be blank")
				os.Exit(1)
			}
			l.Add(task)
		}

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
	case usedFlag["verbose"]:
		listToDoItems(l, true, usedFlag["no-complete"])
	case usedFlag["no-complete"]:
		listToDoItems(l, usedFlag["verbose"], true)
	case len(os.Args) == 1:
		// no arguments provided
		listToDoItems(l, false, false)
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// extracting this logic to make switch statement more readable
// by minimizing indents
func listToDoItems(l *basic_todo_app.List, verbose, noComplete bool) {
	// List current to do items if itemsList is not empty
	if len(*l) <= 0 {
		fmt.Println("You have no to do items")
		os.Exit(0)
	}
	formattedSlice := make([]string, 0)

	if verbose {
		for k, t := range *l {
			prefix := " "
			if t.Done {
				if noComplete {
					continue
				}
				prefix = "X "
				formattedSlice = append(formattedSlice, fmt.Sprintf("%s%d: Completed at %s : %s\n",
					prefix, k+1, t.CompletedAt.Format("Mon, 02 Jan 2006 3:04PM"), t.Task))
				continue
			}
			formattedSlice = append(formattedSlice, fmt.Sprintf("%s%d: Created at %s: %s\n",
				prefix, k+1, t.CreatedAt.Format("Mon, 02 Jan 2006 3:04PM"), t.Task))
		}

		fmt.Print(strings.Join(formattedSlice, ""))
		return
	}
	for k, t := range *l {
		prefix := " "
		if t.Done {
			if noComplete {
				continue
			}
			prefix = "X "
		}

		// Adjust the item number k to print numbers starting from 1
		formattedSlice = append(formattedSlice, fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task))
	}
	fmt.Print(strings.Join(formattedSlice, ""))
}

// getTask function decides where to get the description for a new task
// from: arguments or STDIN
func getTask(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}

	tasks := make([]string, 0)
	s := bufio.NewScanner(r)

	err := s.Err()
	if err != nil {
		return tasks, err
	}

	for s.Scan() {
		tasks = append(tasks, s.Text())
	}

	if len(tasks) == 0 {
		return tasks, fmt.Errorf("task cannot be blank")
	}
	return tasks, nil
}
