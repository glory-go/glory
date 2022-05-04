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
	"github.com/glory-go/glory/log"
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
				log.Errorf("recv error = %s", err)
				return
			}
			fmt.Println(rsp.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchEdit)
}
