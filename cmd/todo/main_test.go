package main_test

import (
	"fmt"
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
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		// Execute command with split string from task variable
		// to simulate multiple arguments
		cmd := exec.Command(cmdPath, strings.Split(task, " ")...)
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		// Execute command with no arguments
		cmd := exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		got := task + "\n"
		if got != string(out) {
			t.Errorf("want %q, got %q", got, string(out))
		}
	})
}
