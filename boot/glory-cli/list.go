package main

import (
	"context"
	"fmt"
)

import (
	"github.com/spf13/cobra"

	"google.golang.org/protobuf/types/known/emptypb"
)

var list = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(defaultDebugAddr)
		rsp, err := debugServiceClient.ListServices(context.Background(), &emptypb.Empty{})
		if err != nil {
			panic(err)
		}
		for _, v := range rsp.ServiceMetadata {
			fmt.Println(v.InterfaceName)
			fmt.Println(v.ImplementationName)
			fmt.Println(v.Methods)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(list)
}
