package main

import (
	"context"
	"fmt"
	"github.com/garethjevans/permissions/asciitree"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func main() {
	c, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		os.Exit(1)
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*c, nil)
	config, err := clientConfig.ClientConfig()

	ctx := context.Background()

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage> permissions <sa>")
	}

	name := os.Args[1]
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		panic(err)
	}

	sa, err := client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	root := asciitree.Tree{}

	clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		if matches(clusterRoleBinding.Subjects, namespace, name) {
			clusterRole, err := client.RbacV1().ClusterRoles().Get(ctx, clusterRoleBinding.RoleRef.Name, metav1.GetOptions{})
			if err != nil {
				panic(err)
			}

			// lets get the permissions
			for _, rule := range clusterRole.Rules {
				for _, resourceName := range rule.Resources {
					root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#ClusterRoleBinding/%s#ClusterRole/%s#%s#%s verbs=%s",
						sa.Name, sa.Namespace, clusterRoleBinding.Name, clusterRole.Name, apiGroup(rule.APIGroups), resourceName, rule.Verbs))
				}
			}
		}
	}

	roleBindings, err := client.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, roleBinding := range roleBindings.Items {
		if matches(roleBinding.Subjects, namespace, name) {
			role, err := client.RbacV1().ClusterRoles().Get(ctx, roleBinding.RoleRef.Name, metav1.GetOptions{})
			if err != nil {
				panic(err)
			}

			// lets get the permissions
			for _, rule := range role.Rules {
				for _, resourceName := range rule.Resources {
					root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)#Role/%s (%s)#%s.%s verbs=%s", sa.Name, sa.Namespace, roleBinding.Name, roleBinding.Namespace, role.Name, role.Namespace, resourceName, apiGroup(rule.APIGroups), rule.Verbs))
				}
			}
		}
	}

	root.Fprint(os.Stdout, true, "")
}

func matches(subjects []v1.Subject, namespace string, name string) bool {
	for _, sub := range subjects {
		if sub.Kind == "ServiceAccount" && sub.Name == name && sub.Namespace == namespace {
			return true
		}
	}
	return false
}

func apiGroup(in []string) string {
	if len(in) == 0 {
		return "empty"
	} else if len(in) == 1 {
		return in[0]
	} else {
		panic("expected length 1")
	}
}
