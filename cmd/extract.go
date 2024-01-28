package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dsnet/compress/bzip2"
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
	Run: func(cmd *cobra.Command, args []string) {
		extractInputFileName := extractInputFileName
		extractOutputDirName := extractOutputDirName
		if extractInputFileName == "" {
			panic("input option is empty")
		}
		if extractOutputDirName == "" {
			panic("output option is empty")
		}

		if err := extractTarNoCompression(extractInputFileName, extractOutputDirName); err != nil {
			panic(err)
		}
		fmt.Printf("extract tar %s %s\n", extractInputFileName, extractOutputDirName)
	},
}

var extractTarGzCmd = &cobra.Command{
	Use: "tar.gz",
	Run: func(cmd *cobra.Command, args []string) {
		extractInputFileName := extractInputFileName
		extractOutputDirName := extractOutputDirName
		if extractInputFileName == "" {
			panic("input option is empty")
		}
		if extractOutputDirName == "" {
			panic("output option is empty")
		}

		if err := extractTarGz(extractInputFileName, extractOutputDirName); err != nil {
			panic(err)
		}
		fmt.Printf("extract tar.gz %s %s\n", extractInputFileName, extractOutputDirName)
	},
}

var extractTarBz2Cmd = &cobra.Command{
	Use: "tar.bz2",
	Run: func(cmd *cobra.Command, args []string) {
		extractInputFileName := extractInputFileName
		extractOutputDirName := extractOutputDirName
		if extractInputFileName == "" {
			panic("input option is empty")
		}
		if extractOutputDirName == "" {
			panic("output option is empty")
		}

		if err := extractTarBz2(extractInputFileName, extractOutputDirName); err != nil {
			panic(err)
		}
		fmt.Printf("extract tar.bz2 %s %s\n", extractInputFileName, extractOutputDirName)
	},
}

func extractZip(extractInputFileName, extractOutputDirName string) error {
	r, err := zip.OpenReader(extractInputFileName)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fmt.Printf("zip file name: %s\n", f.Name)
		if err := copyZip(f, extractOutputDirName); err != nil {
			return err
		}
	}

	return nil
}

func copyZip(f *zip.File, extractOutputDirName string) error {
	dir := filepath.Dir(f.Name)
	if dir != "" && dir != "." {
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

func extractTarNoCompression(extractInputFileName, extractOutputDirName string) error {
	f, err := os.Open(extractInputFileName)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(f)

	return extractTar(extractOutputDirName, tarReader)
}

func extractTarGz(extractInputFileName, extractOutputDirName string) error {
	f, err := os.Open(extractInputFileName)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	return extractTar(extractOutputDirName, tarReader)
}

func extractTarBz2(extractInputFileName, extractOutputDirName string) error {
	f, err := os.Open(extractInputFileName)
	if err != nil {
		return err
	}

	bzr, err := bzip2.NewReader(f, nil)
	if err != nil {
		return err
	}
	defer bzr.Close()

	tarReader := tar.NewReader(bzr)

	return extractTar(extractOutputDirName, tarReader)
}

func extractTar(extractOutputDirName string, tarReader *tar.Reader) error {
	for {
		isEOF, err := copyTar(tarReader, extractOutputDirName)
		if err != nil {
			return err
		}
		if isEOF {
			break
		}
	}

	return nil
}

func copyTar(tarReader *tar.Reader, extractOutputDirName string) (bool, error) {
	hdr, err := tarReader.Next()
	if err != nil {
		if err == io.EOF {
			return true, nil
		}
		return true, err
	}

	tarFileName := hdr.Name
	fmt.Printf("tar file name: %s\n", tarFileName)

	dir := filepath.Dir(tarFileName)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(filepath.Join(extractOutputDirName, dir), 0755); err != nil {
			return true, err
		}
		fmt.Printf("mkdir: %s\n", filepath.Join(extractOutputDirName, dir))
	}

	dst, err := os.Create(filepath.Join(extractOutputDirName, tarFileName))
	if err != nil {
		return true, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, tarReader); err != nil {
		return true, err
	}

	return false, nil
}
