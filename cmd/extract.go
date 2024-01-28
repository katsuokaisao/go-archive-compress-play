package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use: "extract",
}

var extractZipCmd = &cobra.Command{
	Use: "zip",
	Run: func(cmd *cobra.Command, args []string) {
		extractInputFileName := extractInputFileName
		extractOutputDirName := extractOutputDirName
		if extractInputFileName == "" {
			panic("input option is empty")
		}
		if extractOutputDirName == "" {
			panic("output option is empty")
		}

		if err := extractZip(extractInputFileName, extractOutputDirName); err != nil {
			panic(err)
		}
		fmt.Printf("extract zip %s %s\n", extractInputFileName, extractOutputDirName)
	},
}

var extractTarCmd = &cobra.Command{
	Use: "tar",
}

var extractTarGzCmd = &cobra.Command{
	Use: "tar.gz",
}

var extractTarBz2Cmd = &cobra.Command{
	Use: "tar.bz2",
}

func extractZip(extractInputFileName, extractOutputDirName string) error {
	r, err := zip.OpenReader(extractInputFileName)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fmt.Printf("zip file name: %s\n", f.Name)
		if err := copyToZip(f, extractOutputDirName); err != nil {
			return err
		}
	}

	return nil
}

func copyToZip(f *zip.File, extractOutputDirName string) error {
	dir := filepath.Dir(f.Name)
	if dir != "" {
		if err := os.MkdirAll(filepath.Join(extractOutputDirName, dir), 0755); err != nil {
			return err
		}
		fmt.Printf("mkdir: %s\n", filepath.Join(extractOutputDirName, dir))
	}

	dst, err := os.Create(filepath.Join(extractOutputDirName, f.Name))
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
