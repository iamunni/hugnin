/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/iamunni/hugnin/service"
	"github.com/iamunni/hugnin/store"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		noteService := service.NewNoteService(store.NewSQLiteStore())
		err := noteService.View(note)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.PersistentFlags().StringVarP(&note.Value, "note", "n", "", "Search by note")
	viewCmd.PersistentFlags().StringVarP(&note.Tag, "tags", "t", "", "Search by tags")
}
