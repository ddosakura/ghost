package main

import (
	"github.com/ddosakura/ghost/cmd"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ghost-node",
		Short: "Service Node for Ghost-Net",
		Long:  `Ghost-Net is a toolbox of network.`,
	}
)

func main() {
	rootCmd.Version = cmd.Version
	rootCmd.Execute()
}

func init() {

}
