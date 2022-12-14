package main

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/garethjevans/kubectl-permissions/pkg/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-permissions", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdPermissions(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
