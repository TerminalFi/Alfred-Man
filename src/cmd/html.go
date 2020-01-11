package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	htmlPath, err := checkHTML(section, page)
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

func checkHTML(section string, page string) (string, error) {
	path := fmt.Sprintf("%s/%s-%s.html", htmlDir, section, page)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = createHTML(section, page, path)
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

func createHTML(section string, page string, path string) error {
	cmd := exec.Command("/usr/bin/man", "-w", "-s", section, page)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("/usr/bin/groff", "-T", "html", "-mandoc", "-c", strings.TrimSuffix(outb.String(), "\n"))
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()

	if err != nil {
		return errors.New(errb.String())
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.New(errb.String())
	}
	defer f.Close()
	f.Write([]byte(outb.String()))

	return nil
}
