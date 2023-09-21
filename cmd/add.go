/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/iamunni/hugnin/service"
	"github.com/iamunni/hugnin/writer"
	"github.com/spf13/cobra"
)

var note string
var tag string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		value := strings.Join(args, " ")
		noteService := service.NewNoteService(writer.NewSQLiteWriter())
		err := noteService.Add(value, tag)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Tag for the note")
}
