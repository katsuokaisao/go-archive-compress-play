package cmd

import "github.com/spf13/cobra"

var extractCmd = &cobra.Command{
	Use: "extract",
}

var extractZipCmd = &cobra.Command{
	Use: "zip",
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
