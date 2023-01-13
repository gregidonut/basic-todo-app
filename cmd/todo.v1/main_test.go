package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	t.Run("AddNewTask", func(t *testing.T) {
		// Execute command with split string from task1 variable
		// to simulate multiple arguments
		cmd := exec.Command(cmdPath, "-task", task1)
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-task", task2)
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

		want := fmt.Sprintf("%s\n%s\n", task1, task2)
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

	// this should only output one task == task2, since whe 'completed' task1
	t.Run("ListTasksByRunningWithoutFlags", func(t *testing.T) {
		// Execute command with no arguments
		cmd := exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		want := fmt.Sprintf("%s\n", task2)
		if want != string(out) {
			t.Errorf("want %q, got %q", want, string(out))
		}
	})
}
