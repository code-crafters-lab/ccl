package main

import (
	"os"

	"github.com/code-crafters-lab/ccl/cmd"
	"github.com/spf13/cobra"
)

func main() {
	args := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin, args, nil)
	cobra.CheckErr(rootCmd.Execute())
}
