package main

import (
	"fmt"
	"github.com/gregidonut/basic_todo_app"
	"os"
	"strings"
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

	// Decide what to do based on the number of arguments provided
	switch {
	case len(os.Args) == 1:
		// check if there's anything in the List of todo items
		// before looping through each item and then printing it out
		if len(*l) <= 0 {
			fmt.Println("You have no to do items")
			os.Exit(0)
		}
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	// Concatenate all provided arguments with a space and
	// add to the list as an item
	default:
		item := strings.Join(os.Args[1:], " ")
		l.Add(item)
		err := l.Save(todoFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
