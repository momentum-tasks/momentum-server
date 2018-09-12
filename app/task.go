package app

import (
	"time"
)

var taskTableStatement = `CREATE TABLE IF NOT EXISTS tasks (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	owner INT UNSIGNED NOT NULL,
	name TEXT NULL,
	description TEXT NULL,
	due DATETIME,
	priority INT UNSIGNED NOT NULL,
	completed BOOLEAN,
	PRIMARY KEY (id),
	CONSTRAINT UC_Task UNIQUE (owner,priority)
)`

// Task is the basic structure of all tasks
// tasks are associated with a user, and each contain many reports, and other info to define all tasks and events
type Task struct {
	ID          int
	Owner       int
	Name        string
	Description string
	DueDate     time.Time
	Priority    int
	Reports     []Report
	Completed   bool
}

// CreateTask creates a task, stores it in the database, and appends it to the User that was passed in
func CreateTask(user *User, name string, description string, due time.Time, priority int, completed bool) error {
	stmt, err := store.db.Prepare("INSERT INTO tasks(owner, name, description, due, priority, completed) VALUES(?, ?, ?, ?, ?, ?)")
	res, err := stmt.Exec(user.ID, name, description, due, priority, completed)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	user.Tasks = append(user.Tasks, Task{int(id), user.ID, name, description, due, priority, nil, completed})
	return nil
}

// GetTasksByUser finds all tasks associated with a user, and returns them in a sorted order by priority
func GetTasksByUser(user *User) []Task {
	var tasks []Task

	rows, err := store.db.Query("SELECT id, owner, name, description, due, priority, completed FROM tasks WHERE owner = ? ORDER BY priority", user.ID)
	if err != nil {
		mlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Owner, &t.Name, &t.Description, &t.DueDate, &t.Priority, &t.Completed)
		if err != nil {
			mlog.Error(err)
			return tasks
		}
		t.Reports = GetReportsByTask(&t)
		tasks = append(tasks, t)
	}
	return tasks
}

// UpdateName will update the name of a task
func (task *Task) UpdateName(name string) error {
	stmt, err := store.db.Prepare("UPDATE tasks SET name = ? WHERE owner = ? and priority = ?")
	_, err = stmt.Exec(name, task.Owner, task.Priority)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDescription will update the description of a task
func (task *Task) UpdateDescription(desc string) error {
	stmt, err := store.db.Prepare("UPDATE tasks SET description = ? WHERE owner = ? and priority = ?")
	_, err = stmt.Exec(desc, task.Owner, task.Priority)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDue will update the due date of a task
func (task *Task) UpdateDue(due time.Time) error {
	stmt, err := store.db.Prepare("UPDATE tasks SET due = ? WHERE owner = ? and priority = ?")
	_, err = stmt.Exec(due, task.Owner, task.Priority)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePriority inserts the task into a priority, and shuffles all other tasks either up or down to eliminate the gaps
func (task *Task) UpdatePriority(user *User, newPriority int) error {
	var moveDirection int
	if task.Priority > newPriority {
		moveDirection = 1
	} else {
		moveDirection = -1
	}
	stmt, err := store.db.Prepare("UPDATE tasks SET priority = ? WHERE owner = ? and priority = ?")
	_, err = stmt.Exec(0, user.ID, task.Priority)
	if err != nil {
		return err
	}

	for _, t := range user.Tasks {
		if moveDirection > 0 {
			if t.Priority >= newPriority && t.Priority < task.Priority {
				_, err = stmt.Exec(t.Priority+moveDirection, user.ID, t.Priority)
				if err != nil {
					return err
				}
			}
		} else if moveDirection < 0 {
			if t.Priority <= newPriority && t.Priority > task.Priority {
				_, err = stmt.Exec(t.Priority+moveDirection, user.ID, t.Priority)
				if err != nil {
					return err
				}
			}
		}
	}
	_, err = stmt.Exec(newPriority, user.ID, 0)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCompleted will update the completed state of a task
func (task *Task) UpdateCompleted(completed bool) error {
	stmt, err := store.db.Prepare("UPDATE tasks SET completed = ? WHERE owner = ? and priority = ?")
	_, err = stmt.Exec(completed, task.Owner, task.Priority)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes all associated reports, and the task from the database
func (task *Task) Delete() error {
	for _, r := range task.Reports {
		err := r.Delete()
		if err != nil {
			return err
		}
	}
	stmt, err := store.db.Prepare("DELETE FROM tasks WHERE owner = ? and priority = ?")
	_, err = stmt.Exec(task.Owner, task.Priority)
	if err != nil {
		return err
	}
	return nil
}

// DefragTasks works to eliminate gaps in priorities that could be made from a malformed API call, or from deleting a task
func DefragTasks(user *User) error {
	lastPriority := 0
	stmt, err := store.db.Prepare("UPDATE tasks SET priority = ? WHERE owner = ? and priority = ?")
	if err != nil {
		return err
	}

	for _, t := range GetTasksByUser(user) {
		if t.Priority > lastPriority+1 {
			_, err = stmt.Exec(lastPriority+1, user.ID, t.Priority)
			if err != nil {
				return err
			}
			lastPriority++
		} else {
			lastPriority = t.Priority
		}
	}
	return nil
}
