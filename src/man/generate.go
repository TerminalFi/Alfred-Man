package man

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/b4b4r07/go-pipe"
)

func GenerateHTML(section string, page string) (string, error) {
	return generate(section, page, "html")
}

func GeneratePDF(section string, page string) (string, error) {
	return generate(section, page, "pdf")
}

func generate(section string, page string, fileType string) (string, error) {
	path := fmt.Sprintf("data/%s/%s-%s.%s", fileType, section, page, fileType)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = createFile(section, page, path, fileType)
		if err != nil {
			return "", err
		}
		info, _ = os.Stat(path)
	}

	if isOlderThan30Day(info.ModTime()) {
		err = createFile(section, page, path, fileType)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func createFile(section string, page string, path string, fileType string) error {
	if strings.ToLower(fileType) == "pdf" {
		var b bytes.Buffer
		if err := pipe.Command(&b,
			exec.Command("/usr/bin/man", "-t", "-s", section, page),
			exec.Command("/usr/bin/pstopdf", "-i", "-o", path),
		); err != nil {
			return err
		}
	} else if strings.ToLower(fileType) == "html" {
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
	}

	return nil
}

func isOlderThan30Day(t time.Time) bool {
	return time.Now().Sub(t) > (30 * (24 * time.Hour))
}
