package main

import (
	"fmt"
	"log"
)

import (
	"github.com/spf13/cobra"

	"google.golang.org/grpc"
)

import (
	"github.com/glory-go/glory/boot/api/glory/boot"
)

const (
	defaultDebugAddr = "localhost:1999"
)

var rootCmd = &cobra.Command{
	Use: "glory-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func getDebugServiceClent(addr string) boot.DebugServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return boot.NewDebugServiceClient(conn)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}
