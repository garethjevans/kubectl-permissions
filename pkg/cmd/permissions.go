package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/garethjevans/permissions/pkg/asciitree"
	"github.com/kyokomi/emoji/v2"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	permissionsExample = `
	# view the permissions for the specified service account
	%[1]s permissions default
`
	green = ansi.ColorFunc("green")
	red   = ansi.ColorFunc("red")
)

// PermissionsOptions provides information to view permissions
type PermissionsOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

// NewPermissionsOptions provides an instance of PermissionsOptions with default values
func NewPermissionsOptions(streams genericclioptions.IOStreams) *PermissionsOptions {
	return &PermissionsOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

// NewCmdPermissions provides a cobra command wrapping PermissionsOptions
func NewCmdPermissions(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewPermissionsOptions(streams)

	cmd := &cobra.Command{
		Use:          "permissions [name] [flags]",
		Short:        "View the permissions inherited by the specified service account",
		Example:      fmt.Sprintf(permissionsExample, "kubectl"),
		SilenceUsage: true,
		Args:         cobra.ExactValidArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Run lists all available namespaces on a user's KUBECONFIG or updates the
// current context based on a provided namespace.
func (o *PermissionsOptions) Run() error {

	var err error
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	ctx := context.Background()

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage> kubectl permissions <sa>")
	}

	name := os.Args[1]

	namespace := getNamespace(o.configFlags)

	sa, err := client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	root := asciitree.Tree{}

	clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		if matches(clusterRoleBinding.Subjects, namespace, name) {
			clusterRole, err := client.RbacV1().ClusterRoles().Get(ctx, clusterRoleBinding.RoleRef.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Println(red(emoji.Sprint(":no_entry: WARNING")), err)
				root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#ClusterRoleBinding/%s#ClusterRole/%s %s- %s",
					sa.Name,
					sa.Namespace,
					clusterRoleBinding.Name,
					clusterRoleBinding.RoleRef.Name,
					emoji.Sprint(":cross_mark:"),
					red("MISSING!!")))
			} else {
				// lets get the permissions
				for _, rule := range clusterRole.Rules {
					for _, resourceName := range rule.Resources {
						for _, apiGroup := range rule.APIGroups {
							root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#ClusterRoleBinding/%s#ClusterRole/%s#%s#%s verbs=%s %s",
								sa.Name,
								sa.Namespace,
								clusterRoleBinding.Name,
								clusterRole.Name,
								getApiGroup(apiGroup),
								resourceName,
								rule.Verbs,
								green(emoji.Sprint(":check_mark:"))))
						}
					}
				}
			}
		}
	}

	roleBindings, err := client.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, roleBinding := range roleBindings.Items {
		if matches(roleBinding.Subjects, namespace, name) {
			role, err := client.RbacV1().Roles(namespace).Get(ctx, roleBinding.RoleRef.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Println(red(emoji.Sprint(":no_entry: WARNING")), err)
				root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)#Role/%s (%s) %s- %s",
					sa.Name,
					sa.Namespace,
					roleBinding.Name,
					roleBinding.Namespace,
					roleBinding.RoleRef.Name,
					roleBinding.RoleRef.Name,
					emoji.Sprint(":cross_mark:"),
					red("MISSING!!")))
			} else {
				// lets get the permissions
				for _, rule := range role.Rules {
					for _, resourceName := range rule.Resources {
						for _, apiGroup := range rule.APIGroups {
							root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)#Role/%s (%s)#%s#%s verbs=%s %s",
								sa.Name,
								sa.Namespace,
								roleBinding.Name,
								roleBinding.Namespace,
								role.Name,
								role.Namespace,
								getApiGroup(apiGroup),
								resourceName,
								rule.Verbs,
								green(emoji.Sprint(":check_mark:"))))
						}
					}
				}
			}
		}
	}

	root.Fprint(os.Stdout, true, "")

	return nil
}

func matches(subjects []v1.Subject, namespace string, name string) bool {
	for _, sub := range subjects {
		if sub.Kind == "ServiceAccount" && sub.Name == name && sub.Namespace == namespace {
			return true
		}
	}
	return false
}

func getApiGroup(in string) string {
	if len(in) == 0 {
		return "core.k8s.io"
	}
	return in
}
