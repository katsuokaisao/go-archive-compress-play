package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "",
}

var (
	archiveInputDirName      string
	archiveOutputFileName    string
	extractInputFileName     string
	extractOutputDirName     string
	compressInputFileName    string
	compressOutputFileName   string
	decompressInputFileName  string
	decompressOutputFileName string
)

func init() {
	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(compressCmd)
	rootCmd.AddCommand(decompressCmd)

	archiveCmd.AddCommand(archiveZipCmd)
	archiveCmd.AddCommand(archiveTarCmd)
	archiveCmd.AddCommand(archiveTarGzCmd)
	archiveCmd.AddCommand(archiveTarBz2Cmd)

	extractCmd.AddCommand(extractZipCmd)
	extractCmd.AddCommand(extractTarCmd)
	extractCmd.AddCommand(extractTarGzCmd)
	extractCmd.AddCommand(extractTarBz2Cmd)

	compressCmd.AddCommand(compressGzipCmd)
	compressCmd.AddCommand(compressBzip2Cmd)

	decompressCmd.AddCommand(decompressGzipCmd)
	decompressCmd.AddCommand(decompressBzip2Cmd)

	archiveCmd.PersistentFlags().StringVarP(&archiveInputDirName, "input", "i", "sample/archive/input", "input directory")
	archiveCmd.PersistentFlags().StringVarP(&archiveOutputFileName, "output", "o", "", "output file")

	extractCmd.PersistentFlags().StringVarP(&extractInputFileName, "input", "i", "", "input file")
	extractCmd.PersistentFlags().StringVarP(&extractOutputDirName, "output", "o", "sample/extract/output", "output directory")

	compressCmd.PersistentFlags().StringVarP(&compressInputFileName, "input", "i", "", "input file")
	compressCmd.PersistentFlags().StringVarP(&compressOutputFileName, "output", "o", "", "output file")

	decompressCmd.PersistentFlags().StringVarP(&decompressInputFileName, "input", "i", "", "input file")
	decompressCmd.PersistentFlags().StringVarP(&decompressOutputFileName, "output", "o", "", "output file")
}

func Execute() error {
	return rootCmd.Execute()
}
