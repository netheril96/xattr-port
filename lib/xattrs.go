package lib

import (
	"database/sql"
	"fmt"
	"io/fs"         // For fs.DirEntry
	"path/filepath" // For directory walking and relative path calculation

	"github.com/pkg/xattr" // For extended attribute operations (e.g., github.com/pkg/xattr)
	"github.com/rs/zerolog/log"
)

// ExportXattrs recursively walks the given rootDir, reads extended attributes
// for each file, and stores them in the SQLite database.
// It uses the github.com/pkg/xattr library for xattr operations.
func ExportXattrs(rootDir string, db *sql.DB) error {
	// Ensure the table exists
	if err := CreateXattrTable(db); err != nil {
		return fmt.Errorf("failed to create xattr table: %w", err)
	}

	walkFn := func(currentPath string, d fs.DirEntry, errIn error) error {
		if errIn != nil {
			// Error accessing path, e.g., permission issue.
			return fmt.Errorf("error accessing path %s: %w", currentPath, errIn)
		}

		// Calculate relative path to store in the DB
		relativePath, err := filepath.Rel(rootDir, currentPath)
		if err != nil {
			// Should not happen if currentPath is within rootDir
			return fmt.Errorf("failed to get relative path for %s (base: %s): %w", currentPath, rootDir, err)
		}

		// List extended attributes for the file
		attrNames, err := xattr.LList(currentPath)
		if err != nil {
			// This can happen if xattrs are not supported or due to permissions.
			// Depending on requirements, one might log this and continue.
			// xattr.LList returns an empty slice and nil error if no xattrs exist.
			log.Error().Err(err).Str("path", currentPath).Msg("Failed to list xattrs")
			return nil // Continue walking other files/directories
		}

		for _, attrName := range attrNames {
			attrValue, err := xattr.LGet(currentPath, attrName)
			if err != nil {
				log.Error().Err(err).Str("path", currentPath).Str("attrName", attrName).Msg("Failed to get xattr")
				continue // Continue with the next attribute for this file
			}

			// If inserting into DB fails, we still return that error as it's not an xattr.XXX method error.
			if err := InsertXattrRow(db, relativePath, attrName, attrValue); err != nil {
				return fmt.Errorf("failed to insert xattr for %s (name: %s): %w", relativePath, attrName, err)
			}
		}
		return nil
	}

	return filepath.WalkDir(rootDir, walkFn)
}

// ImportXattrs reads extended attributes from the SQLite database and applies them
// to files under the given rootDir.
func ImportXattrs(rootDir string, db *sql.DB) error {
	importFn := func(relativePath string, xattrName string, xattrValue []byte) error {
		// Construct the full path from the root directory and the relative path
		fullPath := filepath.Join(rootDir, relativePath)

		if err := xattr.LSet(fullPath, xattrName, xattrValue); err != nil {
			log.Error().Err(err).Str("path", fullPath).Str("attrName", xattrName).Msg("Failed to set xattr")
			return nil // Continue with the next xattr from the database
		}
		return nil
	}

	if err := IterateXattrRows(db, importFn); err != nil {
		return fmt.Errorf("failed to import xattrs: %w", err)
	}

	return nil
}
