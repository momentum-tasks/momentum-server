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
