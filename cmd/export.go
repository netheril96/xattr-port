/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/netheril96/xattr-port/lib"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all extended attributes under a directory into a database",
	Run: func(cmdCobra *cobra.Command, args []string) {
		db, err := sql.Open("sqlite", globals.DatabasePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open database")
		}
		defer db.Close()

		if err := lib.ExportXattrs(globals.RootDir, db); err != nil {
			log.Error().Err(err).Msg("Failed to export xattrs")
			os.Exit(1)
		}
		log.Info().Msg("Successfully exported extended attributes.")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
