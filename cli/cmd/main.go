package main

import (
	"fmt"
	"os"

	"github.com/vishal-chdhry/policy-reports-extension-api/cli/pkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
