package main

import (
	"context"
	"fmt"
	"log"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/glory-go/glory/debug/api/glory/boot"
)

var watchEdit = &cobra.Command{
	Use: "watchEdit",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(defaultDebugAddr)
		watchEditClient, err := debugServiceClient.WatchEdit(context.Background())
		if err != nil {
			panic(err)
		}
		if err := watchEditClient.Send(&boot.WatchEditRequest{
			InterfaceName:      args[0],
			ImplementationName: args[1],
			Method:             args[2],
			IsParam:            args[3] == "true",
			IsEdit:             false,
		}); err != nil {
			panic(err)
		}
		for {
			rsp, err := watchEditClient.Recv()
			if err != nil {
				log.Printf("recv error = %s\n", err)
				return
			}
			fmt.Println(rsp.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchEdit)
}
