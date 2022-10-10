package main

import (
	"os"

	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"github.com/garethjevans/permissions/pkg/cmd"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-permissions", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdPermissions(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
