package main

import (
	"context"
	"fmt"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/glory-go/glory/boot/api/glory/boot"
)

var watch = &cobra.Command{
	Use: "watch",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(defaultDebugAddr)
		client, err := debugServiceClient.Watch(context.Background(), &boot.WatchRequest{
			InterfaceName:      args[0],
			ImplementationName: args[1],
			Method:             args[2],
			IsParam:            args[3] == "true",
		})
		if err != nil {
			panic(err)
		}
		for {
			msg, _ := client.Recv()
			fmt.Printf(msg.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(watch)
}
