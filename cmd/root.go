/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var globals struct {
	DatabasePath string
	RootDir      string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xattr-port",
	Short: "A simple tool to import and export extended attributes recursively",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&globals.RootDir, "dir", "", "The root dir path")
	rootCmd.PersistentFlags().StringVar(&globals.DatabasePath, "db", "", "The database path")
}
