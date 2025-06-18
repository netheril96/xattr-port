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

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import all extended attributes from a database to a directory",
	Run: func(cmdCobra *cobra.Command, args []string) {
		db, err := sql.Open("sqlite", globals.DatabasePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open database")
		}
		defer db.Close()

		if err := lib.ImportXattrs(globals.RootDir, db); err != nil {
			log.Error().Err(err).Msg("Failed to import xattrs")
			os.Exit(1)
		}
		log.Info().Msg("Successfully imported extended attributes.")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
