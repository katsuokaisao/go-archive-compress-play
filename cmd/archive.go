package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use: "archive",
}

var archiveZipCmd = &cobra.Command{
	Use: "zip",
	Run: func(cmd *cobra.Command, args []string) {
		archiveInputDirName := archiveInputDirName
		archiveOutputFileName := archiveOutputFileName
		if archiveInputDirName == "" {
			panic("input option is empty")
		}
		if archiveOutputFileName == "" {
			panic("output option is empty")
		}

		if err := archiveZip(archiveInputDirName, archiveOutputFileName); err != nil {
			panic(err)
		}
		fmt.Printf("archive zip %s %s\n", archiveInputDirName, archiveOutputFileName)
	},
}

var archiveTarCmd = &cobra.Command{
	Use: "tar",
}

var archiveTarGzCmd = &cobra.Command{
	Use: "tar.gz",
}

var archiveTarBz2Cmd = &cobra.Command{
	Use: "tar.bz2",
}

func archiveZip(archiveInputDirName, archiveOutputFileName string) error {
	dest, err := os.Create(archiveOutputFileName)
	if err != nil {
		return err
	}
	defer dest.Close()

	zipWriter := zip.NewWriter(dest)
	defer zipWriter.Close()

	if err := filepath.Walk(archiveInputDirName, func(path string, info fs.FileInfo, err error) error {
		fmt.Printf("path: %s\n", path)
		if info.IsDir() {
			return nil
		}

		if err := addToZip(path, archiveInputDirName, zipWriter); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func addToZip(srcFileName, archiveInputDirName string, zipWriter *zip.Writer) error {
	src, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	defer src.Close()

	zipFileName, err := filepath.Rel(archiveInputDirName, srcFileName)
	if err != nil {
		return err
	}
	fmt.Printf("zipFileName: %s\n", zipFileName)

	dst, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   zipFileName,
		Method: zip.Deflate,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
