package main

import (
	"log"

	"github.com/xiayinchang/uk8sctl/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
