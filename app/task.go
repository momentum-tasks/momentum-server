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
		mlog.Error(err)
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
