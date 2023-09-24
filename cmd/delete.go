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

var deleteAll bool

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		noteService := service.NewNoteService(store.NewSQLiteStore())
		if deleteAll {
			note.Id = -1
		}
		err := noteService.Delete(note)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.PersistentFlags().Int64VarP(&note.Id, "id", "i", 0, "Search by note Id")
	deleteCmd.PersistentFlags().StringVarP(&note.Tag, "tags", "t", "", "Search by Tags")
	deleteCmd.PersistentFlags().BoolVarP(&deleteAll, "all", "a", false, "Delete All Records")
}
