package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dsnet/compress/bzip2"
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
	Run: func(cmd *cobra.Command, args []string) {
		archiveInputDirName := archiveInputDirName
		archiveOutputFileName := archiveOutputFileName
		if archiveInputDirName == "" {
			panic("input option is empty")
		}
		if archiveOutputFileName == "" {
			panic("output option is empty")
		}

		if err := archiveTarNoCompress(archiveInputDirName, archiveOutputFileName); err != nil {
			panic(err)
		}
		fmt.Printf("archive tar %s %s\n", archiveInputDirName, archiveOutputFileName)
	},
}

var archiveTarGzCmd = &cobra.Command{
	Use: "tar.gz",
	Run: func(cmd *cobra.Command, args []string) {
		archiveInputDirName := archiveInputDirName
		archiveOutputFileName := archiveOutputFileName
		if archiveInputDirName == "" {
			panic("input option is empty")
		}
		if archiveOutputFileName == "" {
			panic("output option is empty")
		}

		if err := archiveTarGz(archiveInputDirName, archiveOutputFileName); err != nil {
			panic(err)
		}
		fmt.Printf("archive tar.gz %s %s\n", archiveInputDirName, archiveOutputFileName)
	},
}

var archiveTarBz2Cmd = &cobra.Command{
	Use: "tar.bz2",
	Run: func(cmd *cobra.Command, args []string) {
		archiveInputDirName := archiveInputDirName
		archiveOutputFileName := archiveOutputFileName
		if archiveInputDirName == "" {
			panic("input option is empty")
		}
		if archiveOutputFileName == "" {
			panic("output option is empty")
		}

		if err := archiveTarBz2(archiveInputDirName, archiveOutputFileName); err != nil {
			panic(err)
		}
		fmt.Printf("archive tar.bz2 %s %s\n", archiveInputDirName, archiveOutputFileName)
	},
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

func archiveTarNoCompress(archiveInputDirName, archiveOutputFileName string) error {
	dest, err := os.Create(archiveOutputFileName)
	if err != nil {
		return err
	}
	defer dest.Close()

	tarWriter := tar.NewWriter(dest)
	defer tarWriter.Close()

	return archiveTar(archiveInputDirName, tarWriter)
}

func archiveTar(archiveInputDirName string, tarWriter *tar.Writer) error {
	if err := filepath.Walk(archiveInputDirName, func(path string, info fs.FileInfo, err error) error {
		fmt.Printf("path: %s\n", path)
		if info.IsDir() {
			return nil
		}

		if err := addToTar(path, archiveInputDirName, tarWriter); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func archiveTarGz(archiveInputDirName, archiveOutputFileName string) error {
	dest, err := os.Create(archiveOutputFileName)
	if err != nil {
		return err
	}
	defer dest.Close()

	gzipWriter := gzip.NewWriter(dest)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return archiveTar(archiveInputDirName, tarWriter)
}

func archiveTarBz2(archiveInputDirName, archiveOutputFileName string) error {
	dest, err := os.Create(archiveOutputFileName)
	if err != nil {
		return err
	}
	defer dest.Close()

	bzip2Writer, err := bzip2.NewWriter(dest, nil)
	if err != nil {
		return err
	}
	defer bzip2Writer.Close()

	tarWriter := tar.NewWriter(dest)
	defer tarWriter.Close()

	return archiveTar(archiveInputDirName, tarWriter)
}

func addToTar(srcFileName, archiveInputDirName string, tarWriter *tar.Writer) error {
	src, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	defer src.Close()

	data := make([]byte, 1024)
	if _, err := src.Read(data); err != nil {
		return err
	}

	tarFileName, err := filepath.Rel(archiveInputDirName, srcFileName)
	if err != nil {
		return err
	}
	fmt.Printf("tarFileName: %s\n", tarFileName)

	hdr := &tar.Header{
		Name: tarFileName,
		Mode: 0600,
		Size: int64(len(data)),
	}

	if err := tarWriter.WriteHeader(hdr); err != nil {
		return err
	}

	if _, err := tarWriter.Write(data); err != nil {
		return err
	}

	return nil
}
