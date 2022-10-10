package cmd

import "k8s.io/cli-runtime/pkg/genericclioptions"

func getNamespace(cf *genericclioptions.ConfigFlags) string {
	if v := *cf.Namespace; v != "" {
		return v
	}
	clientConfig := cf.ToRawKubeConfigLoader()
	defaultNamespace, _, err := clientConfig.Namespace()
	if err != nil {
		defaultNamespace = "default"
	}
	return defaultNamespace
}
