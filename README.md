# Basic Todo App

## Abstract

A basic cli todo app that takes input from several sources
including but not limited to:  
 - Standard Input
 - command line parameters

This app can utilize environment variables to modify how the program
behaves.

## Basic Usage
```

  -add
        Task to be included in the ToDo list
  -complete int
        Item to be completed
  -del int
        Item to be deleted
  -list
        List all tasks
  -no-complete
        prevents showing of completed items
  -verbose
        shows completed or created date
```

## Setting up environment variables
It is currently hardcoded to use the `TODO_FILENAME` environment
variable so you can look up in your respective os how to set this.

## Building binary file
### For windows from linux
`GOOS=windows go build -o path/to/file.exe cmd/todo.v1/main.go`