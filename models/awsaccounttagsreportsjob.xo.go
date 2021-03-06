// Package models contains the types for schema 'trackit'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// AwsAccountTagsReportsJob represents a row from 'trackit.aws_account_tags_reports_job'.
type AwsAccountTagsReportsJob struct {
	ID               int       `json:"id"`               // id
	AwsAccountID     int       `json:"aws_account_id"`   // aws_account_id
	Completed        time.Time `json:"completed"`        // completed
	WorkerID         string    `json:"worker_id"`        // worker_id
	Joberror         string    `json:"jobError"`         // jobError
	Spreadsheeterror string    `json:"spreadsheetError"` // spreadsheetError
	Tagsreporterror  string    `json:"tagsReportError"`  // tagsReportError

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the AwsAccountTagsReportsJob exists in the database.
func (aatrj *AwsAccountTagsReportsJob) Exists() bool {
	return aatrj._exists
}

// Deleted provides information if the AwsAccountTagsReportsJob has been deleted from the database.
func (aatrj *AwsAccountTagsReportsJob) Deleted() bool {
	return aatrj._deleted
}

// Insert inserts the AwsAccountTagsReportsJob to the database.
func (aatrj *AwsAccountTagsReportsJob) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if aatrj._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO trackit.aws_account_tags_reports_job (` +
		`aws_account_id, completed, worker_id, jobError, spreadsheetError, tagsReportError` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.Joberror, aatrj.Spreadsheeterror, aatrj.Tagsreporterror)
	res, err := db.Exec(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.Joberror, aatrj.Spreadsheeterror, aatrj.Tagsreporterror)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	aatrj.ID = int(id)
	aatrj._exists = true

	return nil
}

// Update updates the AwsAccountTagsReportsJob in the database.
func (aatrj *AwsAccountTagsReportsJob) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aatrj._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if aatrj._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE trackit.aws_account_tags_reports_job SET ` +
		`aws_account_id = ?, completed = ?, worker_id = ?, jobError = ?, spreadsheetError = ?, tagsReportError = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.Joberror, aatrj.Spreadsheeterror, aatrj.Tagsreporterror, aatrj.ID)
	_, err = db.Exec(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.Joberror, aatrj.Spreadsheeterror, aatrj.Tagsreporterror, aatrj.ID)
	return err
}

// Save saves the AwsAccountTagsReportsJob to the database.
func (aatrj *AwsAccountTagsReportsJob) Save(db XODB) error {
	if aatrj.Exists() {
		return aatrj.Update(db)
	}

	return aatrj.Insert(db)
}

// Delete deletes the AwsAccountTagsReportsJob from the database.
func (aatrj *AwsAccountTagsReportsJob) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aatrj._exists {
		return nil
	}

	// if deleted, bail
	if aatrj._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM trackit.aws_account_tags_reports_job WHERE id = ?`

	// run query
	XOLog(sqlstr, aatrj.ID)
	_, err = db.Exec(sqlstr, aatrj.ID)
	if err != nil {
		return err
	}

	// set deleted
	aatrj._deleted = true

	return nil
}

// AwsAccount returns the AwsAccount associated with the AwsAccountTagsReportsJob's AwsAccountID (aws_account_id).
//
// Generated from foreign key 'aws_account_tags_reports_job_ibfk_1'.
func (aatrj *AwsAccountTagsReportsJob) AwsAccount(db XODB) (*AwsAccount, error) {
	return AwsAccountByID(db, aatrj.AwsAccountID)
}

// AwsAccountTagsReportsJobByID retrieves a row from 'trackit.aws_account_tags_reports_job' as a AwsAccountTagsReportsJob.
//
// Generated from index 'aws_account_tags_reports_job_id_pkey'.
func AwsAccountTagsReportsJobByID(db XODB, id int) (*AwsAccountTagsReportsJob, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, tagsReportError ` +
		`FROM trackit.aws_account_tags_reports_job ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	aatrj := AwsAccountTagsReportsJob{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&aatrj.ID, &aatrj.AwsAccountID, &aatrj.Completed, &aatrj.WorkerID, &aatrj.Joberror, &aatrj.Spreadsheeterror, &aatrj.Tagsreporterror)
	if err != nil {
		return nil, err
	}

	return &aatrj, nil
}

// AwsAccountTagsReportsJobsByAwsAccountID retrieves a row from 'trackit.aws_account_tags_reports_job' as a AwsAccountTagsReportsJob.
//
// Generated from index 'foreign_aws_account'.
func AwsAccountTagsReportsJobsByAwsAccountID(db XODB, awsAccountID int) ([]*AwsAccountTagsReportsJob, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, tagsReportError ` +
		`FROM trackit.aws_account_tags_reports_job ` +
		`WHERE aws_account_id = ?`

	// run query
	XOLog(sqlstr, awsAccountID)
	q, err := db.Query(sqlstr, awsAccountID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*AwsAccountTagsReportsJob{}
	for q.Next() {
		aatrj := AwsAccountTagsReportsJob{
			_exists: true,
		}

		// scan
		err = q.Scan(&aatrj.ID, &aatrj.AwsAccountID, &aatrj.Completed, &aatrj.WorkerID, &aatrj.Joberror, &aatrj.Spreadsheeterror, &aatrj.Tagsreporterror)
		if err != nil {
			return nil, err
		}

		res = append(res, &aatrj)
	}

	return res, nil
}
