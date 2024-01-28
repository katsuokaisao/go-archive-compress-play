package cmd

import "github.com/spf13/cobra"

var compressCmd = &cobra.Command{
	Use: "compress",
}

var compressGzipCmd = &cobra.Command{
	Use: "gzip",
}

var compressBzip2Cmd = &cobra.Command{
	Use: "bzip2",
}
