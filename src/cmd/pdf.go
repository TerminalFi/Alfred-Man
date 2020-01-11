package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/b4b4r07/go-pipe"
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
	pdfPath, err := checkPDF(section, page)
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

func checkPDF(section string, page string) (string, error) {
	path := fmt.Sprintf("%s/%s-%s.pdf", pdfDir, section, page)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = createPDF(section, page, path)
		if err != nil {
			return "", err
		}
		info, _ = os.Stat(path)
	}

	if isOlderThan30Day(info.ModTime()) {
		err = createPDF(section, page, path)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func createPDF(section string, page string, path string) error {
	var b bytes.Buffer
	if err := pipe.Command(&b,
		exec.Command("/usr/bin/man", "-t", "-s", section, page),
		exec.Command("/usr/bin/pstopdf", "-i", "-o", path),
	); err != nil {
		return err
	}

	return nil
}

func isOlderThan30Day(t time.Time) bool {
	return time.Now().Sub(t) > (30 * (24 * time.Hour))
}
