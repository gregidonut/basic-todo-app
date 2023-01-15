package main_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

// TestMain executes the go build tool to create the binary of the cli
// then runs it with arguments to check if outputs are expected
// then cleans up the files after the function is completed.
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	err := build.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %q: %s", binName, err)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

// TestTodoCLI runs the built binary from TestMain()
// attempts to add a task, then compares the added to task to output of
// running the binary without arguments
func TestTodoCLI(t *testing.T) {
	task1 := "test task number 1"
	task2 := "test task number 2"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("ListTasksWithoutYetSettingUpTheFile", func(t *testing.T) {
		cmd := exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		want := "You have no to do items\n"
		if want != string(out) {
			t.Errorf("want %q, got %q", want, string(out))
		}
	})

	t.Run("AddNewTaskFromArguemnts", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task1)
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		err = cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		want := fmt.Sprintf(" 1: %s\n 2: %s\n", task1, task2)
		if want != string(out) {
			t.Errorf("want %q, got %q", want, string(out))
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	// this should output an X prefix on task number 1
	t.Run("ListTasksByRunningWithoutFlags", func(t *testing.T) {
		// Execute command with no arguments
		cmd := exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		want := fmt.Sprintf("X 1: %s\n 2: %s\n", task1, task2)
		if want != string(out) {
			t.Errorf("want %q, got %q", want, string(out))
		}
	})

	t.Run("RunAppWithEnvVarTODO_FILENAME", func(t *testing.T) {
		const TODO_FILENAME = "new-todo.json"

		err := os.Setenv("TODO_FILENAME", TODO_FILENAME)
		if err != nil {
			t.Fatal(err)
		}

		// testing if running with multiple arguments also works
		flagsPlusArgs := []string{"-add"}
		flagsPlusArgs = append(flagsPlusArgs, strings.Split(task1, " ")...)

		cmd := exec.Command(cmdPath, flagsPlusArgs...)
		err = cmd.Run()
		if err != nil {
			t.Fatal(err)
		}

		// check if the filepath for the json was created
		todoFilePath := filepath.Join(dir, TODO_FILENAME)
		_, err = os.Stat(todoFilePath)
		if errors.Is(err, os.ErrNotExist) {
			t.Fatal("todo file path was not created")
		}

		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		want := fmt.Sprintf(" 1: %s\n", task1)
		if want != string(out) {
			t.Errorf("want %q, got %q", want, string(out))
		}

		// cleanup of json file
		_, err = os.Stat(todoFilePath)
		if err == nil {
			os.Remove(todoFilePath)
		}
	})
}
