package app

import (
	"fmt"
	"github.com/duxweb/go-fast/global"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

func Command(command *cobra.Command) {
	version := &cobra.Command{
		Use:   "version",
		Short: "View the version number",
		Run: func(cmd *cobra.Command, args []string) {
			color.Println(fmt.Sprintf("⇨ <red>%s</>", global.Version))
		},
	}

	appList := &cobra.Command{
		Use:   "app:list",
		Short: "Viewing the application list",
		Run: func(cmd *cobra.Command, args []string) {

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Name", "Title", "Desc"})
			rows := []table.Row{}
			for _, config := range List {
				rows = append(rows, table.Row{config.Name, config.Title, config.Desc})
			}
			t.AppendRows(rows)
			t.Render()

		},
	}

	command.AddCommand(version, appList)
}
