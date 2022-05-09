package main

import (
	"context"
	"fmt"
)

import (
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

import (
	"github.com/glory-go/glory/debug/api/glory/boot"
)

var watch = &cobra.Command{
	Use: "watch",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(defaultDebugAddr)
		client, err := debugServiceClient.Watch(context.Background(), &boot.WatchRequest{
			InterfaceName:      args[0],
			ImplementationName: args[1],
			Method:             args[2],
			Input:              true,
			Output:             true,
		})
		if err != nil {
			panic(err)
		}
		for {
			msg, _ := client.Recv()
			fmt.Println()
			onToPrint := "Call"
			paramOrResponse := "Param"
			if !msg.IsParam {
				onToPrint = "Response"
				paramOrResponse = "Response"
			}
			color.Red("========== On %s ==========\n", onToPrint)
			color.Red("%s.(%s).%s()", msg.InterfaceName, msg.ImplementationName, msg.MethodName)
			for index, p := range msg.GetParams() {
				color.Cyan("%s %d: %s", paramOrResponse, index+1, p)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(watch)
}
