package cmd

import (
	"alfredman/man"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const htmlDir = "./data/html"

var htmlCmd = &cobra.Command{
	Use:    "html",
	Short:  "Convert MAN page to HTML",
	PreRun: validateHTMLFolder,
	RunE:   htmlRun,
}

func htmlRun(cmd *cobra.Command, args []string) error {
	page, _ := cmd.Flags().GetString("page")
	section, _ := cmd.Flags().GetString("section")
	htmlPath, err := man.GenerateHTML(section, page)
	if err != nil {
		return err
	}
	fmt.Print(htmlPath)
	return nil
}

func init() {

	htmlCmd.Flags().StringP("section", "s", "", "Section for the MAN page")
	htmlCmd.Flags().StringP("page", "p", "", "MAN page to retrieve")

	rootCmd.AddCommand(htmlCmd)
}

func validateHTMLFolder(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(htmlDir); os.IsNotExist(err) {
		os.MkdirAll(htmlDir, os.ModePerm)
	}
}
