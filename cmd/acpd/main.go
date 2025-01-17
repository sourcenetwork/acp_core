package main

import (
	"log"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.Execute()
}

var (
	grpcAddress string
)

var rootCmd = &cobra.Command{
	Use:   "acpd [flags]",
	Short: "acpd starts a grpc server for acpCore",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := runtime.NewRuntimeManager()
		if err != nil {
			log.Fatal("Failed to create acpCore runtime:", err)
		}

		svr := server.NewServer(grpcAddress)
		err = svr.Init(manager)
		if err != nil {
			log.Fatal("Failed to initialize grpc server:", err)
		}

		log.Printf("Serving gRPC service on http://%s", grpcAddress)
		svr.Run()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&grpcAddress, "address", "a", "0.0.0.0:9090", "GRPC server listener address")
}
