package cmd

import (
	"alfredman/man"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const pdfDir = "./data/pdf"

var pdfCmd = &cobra.Command{
	Use:    "pdf",
	Short:  "Convert MAN page to PDF",
	PreRun: validatePDFFolder,
	RunE:   pdfRun,
}

func pdfRun(cmd *cobra.Command, args []string) error {
	page, _ := cmd.Flags().GetString("page")
	section, _ := cmd.Flags().GetString("section")
	pdfPath, err := man.GeneratePDF(section, page)
	if err != nil {
		return err
	}
	fmt.Print(pdfPath)
	return nil
}

func init() {

	pdfCmd.Flags().StringP("section", "s", "", "Section for the MAN page")
	pdfCmd.Flags().StringP("page", "p", "", "MAN page to retrieve")

	rootCmd.AddCommand(pdfCmd)
}

func validatePDFFolder(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(pdfDir); os.IsNotExist(err) {
		os.MkdirAll(pdfDir, os.ModePerm)
	}
}
