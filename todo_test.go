package basic_todo_app_test

import (
	"github.com/gregidonut/basic_todo_app"
	"os"
	"testing"
)

func TestList_Add(t *testing.T) {
	type args struct {
		task string
	}
	tests := []struct {
		name string
		l    basic_todo_app.List
		args args
	}{
		{
			name: "initial",
			l:    basic_todo_app.List{},
			args: args{
				task: "New Task",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Add(tt.args.task)
			if tt.l[0].Task != tt.args.task {
				t.Errorf("want %q, got %q", tt.args.task, tt.l[0].Task)
			}
		})
	}
}

func TestList_Complete(t *testing.T) {
	type args struct {
		i        int
		taskName string
	}
	tests := []struct {
		name string
		l    basic_todo_app.List
		args args
		//wantErr bool
	}{
		{
			name: "initial",
			l:    basic_todo_app.List{},
			args: args{
				i:        1,
				taskName: "New Task",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Add(tt.args.taskName)

			if tt.l[0].Task != tt.args.taskName {
				t.Errorf("want %q, got %q", tt.args.taskName, tt.l[0].Task)
			}

			if tt.l[0].Done {
				t.Errorf("New task should not be completed but it is")
			}

			tt.l.Complete(tt.args.i)
			if !tt.l[0].Done {
				t.Errorf("New task should be completed but is not")
			}
		})
	}
}

func TestList_Delete(t *testing.T) {
	type args struct {
		i         int
		taskNames []string
	}
	tests := []struct {
		name    string
		l       basic_todo_app.List
		args    args
		wantErr bool
	}{
		{
			name: "initial",
			l:    basic_todo_app.List{},
			args: args{
				i: 2,
				taskNames: []string{
					"New Task 1",
					"New Task 2",
					"New Task 3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.taskNames {
				tt.l.Add(v)
			}

			if tt.l[0].Task != tt.args.taskNames[0] {
				t.Errorf("want %q, got %q", tt.args.taskNames[0], tt.l[0].Task)
			}

			tt.l.Delete(tt.args.i)
			if len(tt.l) != 2 {
				t.Errorf("want List length %d, got %d", tt.args.i, len(tt.l))
			}

			if tt.l[1].Task != tt.args.taskNames[2] {
				t.Errorf("Delete() did not delete the first element, want %q, got %q", tt.args.taskNames[2], tt.l[1].Task)
			}

		})
	}
}

func TestList_SaveGet(t *testing.T) {
	type args struct {
		filename string
		taskName string
	}
	tests := []struct {
		name    string
		l1      basic_todo_app.List
		l2      basic_todo_app.List
		args    args
		wantErr bool
	}{
		{
			name: "initial",
			l1:   basic_todo_app.List{},
			l2:   basic_todo_app.List{},
			args: args{
				taskName: "New Task",
			},
		},
	}
	for _, tt := range tests {
		// here we create 2 Lists save one on a temp file then load it on the second one
		// with the get function. then we check if their contents if both Lists are the same.
		t.Run(tt.name, func(t *testing.T) {
			tt.l1.Add(tt.args.taskName)
			// check to see if first task in List 1 is equal to the task name as a sanity check
			if tt.l1[0].Task != tt.args.taskName {
				t.Errorf("Expected %q, got %q instead.", tt.args.taskName, tt.l1[0].Task)
			}
			// check if there is an error creating the tmp file
			tf, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatalf("Error creating temp file: %s", err)
			}
			defer os.Remove(tf.Name())
			// check if Save() returns an error
			// which should occur if either json.Marshal() has an error or os.WriteFile() has an error
			err = tt.l1.Save(tf.Name())
			if err != nil {
				t.Fatalf("Error saving list to file: %s", err)
			}
			// check if Get() returns an error
			// which should occur if json.UnMarshal() has a error
			// or os.ReadFile() returns an error other than the file not existing
			err = tt.l2.Get(tf.Name())
			if err != nil {
				t.Fatalf("Error getting list from file: %s", err)
			}
			// check if first task in List 1 matches the first task in list 2
			if tt.l1[0].Task != tt.l2[0].Task {
				t.Errorf("Task %q shold match %q task.", tt.l1[0].Task, tt.l2[0].Task)
			}

		})
	}
}
