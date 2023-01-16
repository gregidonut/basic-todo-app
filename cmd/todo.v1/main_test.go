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
	"time"
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
	// empty strings to hold formatted time info
	task2CreationTime := ""
	task1DoneTimeTime := ""

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
		task2CreationTime = time.Now().Format("Mon, 02 Jan 2006 3:04PM")

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

	t.Run("AddBlankTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		io.WriteString(cmdStdIn, "")
		cmdStdIn.Close()

		out, err := cmd.CombinedOutput()
		want := "task cannot be blank\n"
		if want != string(out) {
			t.Errorf("want %q, got %q", want, string(out))
		}
		if err == nil {
			t.Errorf("expected error but didn't get one")
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
		task1DoneTimeTime = time.Now().Format("Mon, 02 Jan 2006 3:04PM")

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

	t.Run("ListVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("X 1: Completed at %s : %s\n 2: Created at %s: %s\n",
			task1DoneTimeTime, task1, task2CreationTime, task2)
		if want != string(out) {
			t.Errorf("\n\t\twant\t %q,\n\t\tgot\t\t %q", want, string(out))
		}
	})

	t.Run("VerboseWithoutList", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}
		want := fmt.Sprintf("X 1: Completed at %s : %s\n 2: Created at %s: %s\n",
			task1DoneTimeTime, task1, task2CreationTime, task2)
		if want != string(out) {
			t.Errorf("\n\t\twant\t %q,\n\t\tgot\t\t %q", want, string(out))
		}
	})

	t.Run("NoCompletedTasksByItself", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-no-complete")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}
		want := fmt.Sprintf(" 2: %s\n", task2)
		if want != string(out) {
			t.Errorf("\n\t\twant\t %q,\n\t\tgot\t\t %q", want, string(out))
		}
	})

	t.Run("NoCompletedTasksWithList", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-no-complete")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}
		want := fmt.Sprintf(" 2: %s\n", task2)
		if want != string(out) {
			t.Errorf("\n\t\twant\t %q,\n\t\tgot\t\t %q", want, string(out))
		}
	})

	t.Run("NoCompletedTasksWithVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-verbose", "-no-complete")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}
		want := fmt.Sprintf(" 2: Created at %s: %s\n",
			task2CreationTime, task2)
		if want != string(out) {
			t.Errorf("\n\t\twant\t %q,\n\t\tgot\t\t %q", want, string(out))
		}
	})

	t.Run("NoCompletedTasksWithListAndVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-verbose", "-no-complete")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}
		want := fmt.Sprintf(" 2: Created at %s: %s\n",
			task2CreationTime, task2)
		if want != string(out) {
			t.Errorf("\n\t\twant\t %q,\n\t\tgot\t\t %q", want, string(out))
		}
	})
	// at this point there should be two items in the list
	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-del", "1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}

		want := fmt.Sprintf("Deleted task number 1: %s\n", task1)
		if want != string(out) {
			t.Errorf("want %q, got %q", want, string(out))
		}
	})

	// at this point there should be one item in the list
	t.Run("AddNewTaskFromMultiLineSTDIN", func(t *testing.T) {
		multilineInput := "task number 3\ntask number 4\ntask number 5\n"

		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		io.WriteString(cmdStdIn, multilineInput)
		cmdStdIn.Close()

		err = cmd.Run()
		if err != nil {
			t.Fatal(err)
		}

		// listing items here to check should be task number 2 through 3
		cmd = exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		formattedMultlineInputSlice := strings.Split(multilineInput, "\n")
		var formattedMultlineInputString string
		for i, v := range formattedMultlineInputSlice[:len(formattedMultlineInputSlice)-1] {
			formattedMultlineInputString += fmt.Sprintf(" %d: %s\n", i+2, v)
		}
		want := fmt.Sprintf(" 1: %s\n%s", task2, formattedMultlineInputString)
		if want != string(out) {
			t.Errorf("\nwant \t\t%q, \ngot \t\t%q", want, string(out))
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
