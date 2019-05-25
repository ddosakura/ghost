package main

import (
	"github.com/armon/go-socks5"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ghost-master",
		Short: "Master-Node for Ghost-Net",
		Long:  `Ghost-Net is a toolbox of network.`,
		Run: func(c *cobra.Command, args []string) {
			conf := &socks5.Config{}
			server, err := socks5.New(conf)
			if err != nil {
				panic(err)
			}

			if err := server.ListenAndServe("tcp", addr); err != nil {
				panic(err)
			}
		},
	}

	addr string
)

func main() {
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "127.0.0.1:4405", "addr of socks")
	rootCmd.Execute()
}
