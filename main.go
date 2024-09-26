package main

import (
	"fmt"

	"github.com/punmin/cloud-bill/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
