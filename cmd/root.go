package cmd

import (
	"log"
	"os"

	"github.com/curefatih/gusr/pkg"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gusr",
	Short: "A tool for managing Git users",
}

func Execute() {
	if err := pkg.CreateConfigFileIfNotExist(); err != nil {
		log.Fatal(err)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
