package lib

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

const tableName = "xattrs"

const createTableSQL = `
CREATE TABLE IF NOT EXISTS %s (
    relative_path TEXT NOT NULL,
    xattr_name TEXT NOT NULL,
    xattr_value BLOB NOT NULL,
    PRIMARY KEY (relative_path, xattr_name)
);`

// CreateXattrTable creates the SQLite table for storing extended attributes if it doesn't already exist.
func CreateXattrTable(db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf(createTableSQL, tableName))
	if err != nil {
		return fmt.Errorf("failed to prepare create table statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute create table statement: %w", err)
	}
	return nil
}

// InsertXattrRow inserts a new row into the xattrs table.
func InsertXattrRow(db *sql.DB, relativePath string, xattrName string, xattrValue []byte) error {
	insertSQL := fmt.Sprintf("INSERT INTO %s (relative_path, xattr_name, xattr_value) VALUES (?, ?, ?)", tableName)
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(relativePath, xattrName, xattrValue)
	if err != nil {
		return fmt.Errorf("failed to execute insert statement: %w", err)
	}
	return nil
}

// IterateXattrRows iterates over all rows in the xattrs table and calls the provided callback function for each row.
func IterateXattrRows(db *sql.DB, callback func(relativePath string, xattrName string, xattrValue []byte) error) error {
	querySQL := fmt.Sprintf("SELECT relative_path, xattr_name, xattr_value FROM %s", tableName)
	rows, err := db.Query(querySQL)
	if err != nil {
		return fmt.Errorf("failed to query rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var relativePath, xattrName string
		var xattrValue []byte
		if err := rows.Scan(&relativePath, &xattrName, &xattrValue); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		if err := callback(relativePath, xattrName, xattrValue); err != nil {
			return fmt.Errorf("callback error: %w", err)
		}
	}
	return rows.Err() // Check for errors during iteration
}
