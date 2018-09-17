package app

var reportTableStatement = `CREATE TABLE IF NOT EXISTS reports (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	taskid INT UNSIGNED NOT NULL,
	sequence INT UNSIGNED NOT NULL,
	description TEXT NOT NULL,
	PRIMARY KEY (id),
	CONSTRAINT UC_Report UNIQUE (taskid,sequence)
)`

// Report is the basic structure of the Report object
// reports are to be attached to a single task, and ordered in a sequence
type Report struct {
	ID          int
	TaskID      int
	Sequence    int
	Description string
}

// CreateReport creates a report and attaches it to a task
func CreateReport(task *Task, sequence int, description string) error {
	stmt, err := store.db.Prepare("INSERT INTO reports(taskid, sequence, description) VALUES(?, ?, ?)")
	res, err := stmt.Exec(task.ID, sequence, description)
	if err != nil {
		mlog.Error(err)
		return err
	}
	id, _ := res.LastInsertId()
	task.Reports = append(task.Reports, Report{int(id), task.ID, sequence, description})
	return nil
}

// GetReportsByTask finds all reports associated with a task, and returns them in their correct sequenced order
func GetReportsByTask(task *Task) []Report {
	var reports []Report

	rows, err := store.db.Query("SELECT id, taskid, sequence, description FROM reports WHERE taskid = ? ORDER BY sequence", task.ID)
	if err != nil {
		mlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r Report
		err := rows.Scan(&r.ID, &r.TaskID, &r.Sequence, &r.Description)
		if err != nil {
			mlog.Error(err)
			return reports
		}
		reports = append(reports, r)
	}
	return reports
}

// UpdateDescription will update the description of the report
func (report *Report) UpdateDescription(desc string) error {
	stmt, err := store.db.Prepare("UPDATE reports SET description = ? WHERE taskid = ? and sequence = ?")
	_, err = stmt.Exec(desc, report.TaskID, report.Sequence)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSequence inserts the report into a sequence, and shuffles all other reports either up or down to eliminate the gaps
func (report *Report) UpdateSequence(task *Task, newSequence int) error {
	var moveDirection int
	if report.Sequence > newSequence {
		moveDirection = 1
	} else {
		moveDirection = -1
	}
	stmt, err := store.db.Prepare("UPDATE tasks SET sequence = ? WHERE task = ? and sequence = ?")
	_, err = stmt.Exec(0, task.ID, report.Sequence)
	if err != nil {
		return err
	}

	for _, r := range task.Reports {
		if moveDirection > 0 {
			if r.Sequence >= newSequence && r.Sequence < task.Priority {
				_, err = stmt.Exec(r.Sequence+moveDirection, task.ID, r.Sequence)
				if err != nil {
					return err
				}
			}
		} else if moveDirection < 0 {
			if r.Sequence <= newSequence && r.Sequence > task.Priority {
				_, err = stmt.Exec(r.Sequence+moveDirection, task.ID, r.Sequence)
				if err != nil {
					return err
				}
			}
		}
	}
	_, err = stmt.Exec(newSequence, task.ID, 0)
	if err != nil {
		return err
	}
	err = DefragReports(task)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the report from the database
func (report *Report) Delete() error {
	stmt, err := store.db.Prepare("DELETE FROM reports WHERE taskid = ? and sequence = ?")
	_, err = stmt.Exec(report.TaskID, report.Sequence)
	if err != nil {
		mlog.Error(err)
		return err
	}
	return nil
}

// DefragReports works to eliminate gaps in sequence that could be made from a malformed API call, or from deleting a report
func DefragReports(task *Task) error {
	lastSequence := 0
	stmt, err := store.db.Prepare("UPDATE reports SET sequence = ? WHERE taskid = ? and sequence = ?")
	if err != nil {
		return err
	}

	for _, r := range GetReportsByTask(task) {
		if r.Sequence > lastSequence+1 {
			_, err = stmt.Exec(lastSequence+1, task.ID, r.Sequence)
			if err != nil {
				return err
			}
			lastSequence++
		} else {
			lastSequence = r.Sequence
		}
	}
	return nil
}
