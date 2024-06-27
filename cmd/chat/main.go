package main

import (
	"chat-server/internal/chat"
	"fmt"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {
	//var configPath string
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "start chat server",
		Example: `./chat start`,
		Run: func(cmd *cobra.Command, args []string) {
			server := chat.NewServer()
			err := server.Run()
			if err != nil {
				fmt.Println("start server error: ", err)
			}
		},
	}

	//cmd.PersistentFlags().StringVarP(&configPath, "configPath", "c", "", "configuration file path for chat server")
	return cmd
}

func main() {
	root := newServerCommand()
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
