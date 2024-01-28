package cmd

import "github.com/spf13/cobra"

var decompressCmd = &cobra.Command{
	Use: "decompress",
}

var decompressGzipCmd = &cobra.Command{
	Use: "gzip",
}

var decompressBzip2Cmd = &cobra.Command{
	Use: "bzip2",
}
