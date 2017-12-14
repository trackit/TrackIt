// Package models contains the types for schema 'trackit'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// ChargifyConfiguration represents a row from 'trackit.chargify_configuration'.
type ChargifyConfiguration struct {
	ID                  int    `json:"id"`                    // id
	AwsAccountID        int    `json:"aws_account_id"`        // aws_account_id
	BillingBucket       string `json:"billing_bucket"`        // billing_bucket
	Subdomain           string `json:"subdomain"`             // subdomain
	APIKey              []byte `json:"api_key"`               // api_key
	ChargifyTagKey      string `json:"chargify_tag_key"`      // chargify_tag_key
	StandardStorageID   int    `json:"standard_storage_id"`   // standard_storage_id
	InfrequentStorageID int    `json:"infrequent_storage_id"` // infrequent_storage_id
	GlacierStorageID    int    `json:"glacier_storage_id"`    // glacier_storage_id
	BandwidthID         int    `json:"bandwidth_id"`          // bandwidth_id
	RequestID           int    `json:"request_id"`            // request_id

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the ChargifyConfiguration exists in the database.
func (cc *ChargifyConfiguration) Exists() bool {
	return cc._exists
}

// Deleted provides information if the ChargifyConfiguration has been deleted from the database.
func (cc *ChargifyConfiguration) Deleted() bool {
	return cc._deleted
}

// Insert inserts the ChargifyConfiguration to the database.
func (cc *ChargifyConfiguration) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cc._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO trackit.chargify_configuration (` +
		`aws_account_id, billing_bucket, subdomain, api_key, chargify_tag_key, standard_storage_id, infrequent_storage_id, glacier_storage_id, bandwidth_id, request_id` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cc.AwsAccountID, cc.BillingBucket, cc.Subdomain, cc.APIKey, cc.ChargifyTagKey, cc.StandardStorageID, cc.InfrequentStorageID, cc.GlacierStorageID, cc.BandwidthID, cc.RequestID)
	res, err := db.Exec(sqlstr, cc.AwsAccountID, cc.BillingBucket, cc.Subdomain, cc.APIKey, cc.ChargifyTagKey, cc.StandardStorageID, cc.InfrequentStorageID, cc.GlacierStorageID, cc.BandwidthID, cc.RequestID)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cc.ID = int(id)
	cc._exists = true

	return nil
}

// Update updates the ChargifyConfiguration in the database.
func (cc *ChargifyConfiguration) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cc._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cc._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE trackit.chargify_configuration SET ` +
		`aws_account_id = ?, billing_bucket = ?, subdomain = ?, api_key = ?, chargify_tag_key = ?, standard_storage_id = ?, infrequent_storage_id = ?, glacier_storage_id = ?, bandwidth_id = ?, request_id = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, cc.AwsAccountID, cc.BillingBucket, cc.Subdomain, cc.APIKey, cc.ChargifyTagKey, cc.StandardStorageID, cc.InfrequentStorageID, cc.GlacierStorageID, cc.BandwidthID, cc.RequestID, cc.ID)
	_, err = db.Exec(sqlstr, cc.AwsAccountID, cc.BillingBucket, cc.Subdomain, cc.APIKey, cc.ChargifyTagKey, cc.StandardStorageID, cc.InfrequentStorageID, cc.GlacierStorageID, cc.BandwidthID, cc.RequestID, cc.ID)
	return err
}

// Save saves the ChargifyConfiguration to the database.
func (cc *ChargifyConfiguration) Save(db XODB) error {
	if cc.Exists() {
		return cc.Update(db)
	}

	return cc.Insert(db)
}

// Delete deletes the ChargifyConfiguration from the database.
func (cc *ChargifyConfiguration) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cc._exists {
		return nil
	}

	// if deleted, bail
	if cc._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM trackit.chargify_configuration WHERE id = ?`

	// run query
	XOLog(sqlstr, cc.ID)
	_, err = db.Exec(sqlstr, cc.ID)
	if err != nil {
		return err
	}

	// set deleted
	cc._deleted = true

	return nil
}

// AwsAccount returns the AwsAccount associated with the ChargifyConfiguration's AwsAccountID (aws_account_id).
//
// Generated from foreign key 'chargify_configuration_ibfk_1'.
func (cc *ChargifyConfiguration) AwsAccount(db XODB) (*AwsAccount, error) {
	return AwsAccountByID(db, cc.AwsAccountID)
}

// ChargifyConfigurationByAwsAccountID retrieves a row from 'trackit.chargify_configuration' as a ChargifyConfiguration.
//
// Generated from index 'aws_account_id'.
func ChargifyConfigurationByAwsAccountID(db XODB, awsAccountID int) (*ChargifyConfiguration, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, billing_bucket, subdomain, api_key, chargify_tag_key, standard_storage_id, infrequent_storage_id, glacier_storage_id, bandwidth_id, request_id ` +
		`FROM trackit.chargify_configuration ` +
		`WHERE aws_account_id = ?`

	// run query
	XOLog(sqlstr, awsAccountID)
	cc := ChargifyConfiguration{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, awsAccountID).Scan(&cc.ID, &cc.AwsAccountID, &cc.BillingBucket, &cc.Subdomain, &cc.APIKey, &cc.ChargifyTagKey, &cc.StandardStorageID, &cc.InfrequentStorageID, &cc.GlacierStorageID, &cc.BandwidthID, &cc.RequestID)
	if err != nil {
		return nil, err
	}

	return &cc, nil
}

// ChargifyConfigurationByID retrieves a row from 'trackit.chargify_configuration' as a ChargifyConfiguration.
//
// Generated from index 'chargify_configuration_id_pkey'.
func ChargifyConfigurationByID(db XODB, id int) (*ChargifyConfiguration, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, billing_bucket, subdomain, api_key, chargify_tag_key, standard_storage_id, infrequent_storage_id, glacier_storage_id, bandwidth_id, request_id ` +
		`FROM trackit.chargify_configuration ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	cc := ChargifyConfiguration{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cc.ID, &cc.AwsAccountID, &cc.BillingBucket, &cc.Subdomain, &cc.APIKey, &cc.ChargifyTagKey, &cc.StandardStorageID, &cc.InfrequentStorageID, &cc.GlacierStorageID, &cc.BandwidthID, &cc.RequestID)
	if err != nil {
		return nil, err
	}

	return &cc, nil
}