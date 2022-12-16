package cmd

import (
	"context"
	"fmt"
	"github.com/garethjevans/kubectl-permissions/pkg/asciitree"
	"github.com/garethjevans/kubectl-permissions/pkg/roles"
	"github.com/garethjevans/kubectl-permissions/pkg/version"
	"github.com/kyokomi/emoji/v2"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
)

const (
	CROSS_MARK = ":cross_mark:"
	CHECK_MARK = ":check_mark:"
	NO_ENTRY   = ":no_entry:"
)

const binaryName = "kubectl"

var (
	permissionsExample = `
	# view the permissions for the 'default' service account
	%[1]s permissions default

	# view the permissions for the 'sa' service account in the namespace 'test'
	%[1]s permissions sa -n tests`

	noColor = os.Getenv("NO_COLOR") == "true"
)

// PermissionsOptions provides information to view permissions
type PermissionsOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	Version bool

	Cmd  *cobra.Command
	Args []string
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
		Use:          fmt.Sprintf("%s permissions [service-account] [flags]", binaryName),
		Short:        "View the permissions inherited by the specified service account",
		Example:      fmt.Sprintf(permissionsExample, binaryName),
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			o.Args = args
			o.Cmd = c
			return o.Run()
		},
	}

	cmd.Flags().BoolVarP(&o.Version, "version", "v", false, "Display the version of the permissions plugin")

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Run lists all available namespaces on a user's KUBECONFIG or updates the
// current context based on a provided namespace.
func (o *PermissionsOptions) Run() error {
	if o.Version {
		fmt.Println(version.Version)
		return nil
	}

	if len(o.Args) != 1 {
		fmt.Println("Usage:", o.Cmd.Use)
		return fmt.Errorf("error: accepts 1 arg(s), received %d", len(o.Args))
	}

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

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return err
	}

	r, err := roles.DiscoverRolesAndPermissions(discoveryClient)
	if err != nil {
		return err
	}

	// there must only be one arg
	name := o.Args[0]

	namespace := getNamespace(o.configFlags)

	sa, err := client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	root := asciitree.Tree{}

	root.Add(fmt.Sprintf("ServiceAccount/%s (%s)",
		sa.Name,
		sa.Namespace))

	clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		if matches(clusterRoleBinding.Subjects, namespace, name) {
			root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#ClusterRoleBinding/%s",
				sa.Name,
				sa.Namespace,
				clusterRoleBinding.Name))

			clusterRole, err := client.RbacV1().ClusterRoles().Get(ctx, clusterRoleBinding.RoleRef.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), err)
				root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#ClusterRoleBinding/%s#ClusterRole/%s %s- %s",
					sa.Name,
					sa.Namespace,
					clusterRoleBinding.Name,
					clusterRoleBinding.RoleRef.Name,
					getEmoji(CROSS_MARK),
					red("MISSING!!")))
			} else {
				// lets get the permissions
				for _, rule := range clusterRole.Rules {
					for _, resourceName := range rule.Resources {
						for _, apiGroup := range rule.APIGroups {
							mark := green(getEmoji(CHECK_MARK))
							message := ""
							availableApiGroup, ok := r[apiGroup]
							if !ok {
								fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), "API Group", apiGroup, "does not exist")
								mark = red(getEmoji(CROSS_MARK))
								message = fmt.Sprintf(" (API Group '%s' does not exist)", apiGroup)
							} else {
								verbs, ok := availableApiGroup[resourceName]
								if !ok {
									fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), "Resource", resourceName, "does not exist")
									mark = red(getEmoji(CROSS_MARK))
									message = fmt.Sprintf(" (Resource '%s' does not exist)", resourceName)
								} else {
									verbMessage, ok := validateVerbs(rule.Verbs, verbs)
									if !ok {
										mark = red(getEmoji(CROSS_MARK))
										message = verbMessage
									}
								}
							}
							root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#ClusterRoleBinding/%s#ClusterRole/%s#%s#%s verbs=%s %s%s",
								sa.Name,
								sa.Namespace,
								clusterRoleBinding.Name,
								clusterRole.Name,
								getApiGroup(apiGroup),
								resourceName,
								rule.Verbs,
								mark,
								message))
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
			root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)",
				sa.Name,
				sa.Namespace,
				roleBinding.Name,
				roleBinding.Namespace))

			if roleBinding.RoleRef.Kind == "Role" {
				role, err := client.RbacV1().Roles(namespace).Get(ctx, roleBinding.RoleRef.Name, metav1.GetOptions{})
				if err != nil {
					fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), err)
					root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)#Role/%s (%s) %s- %s",
						sa.Name,
						sa.Namespace,
						roleBinding.Name,
						roleBinding.Namespace,
						roleBinding.RoleRef.Name,
						roleBinding.RoleRef.Name,
						getEmoji(CROSS_MARK),
						red("MISSING!!")))
				} else {
					// lets get the permissions
					for _, rule := range role.Rules {
						for _, resourceName := range rule.Resources {
							for _, apiGroup := range rule.APIGroups {
								mark := green(getEmoji(CHECK_MARK))
								message := ""
								availableApiGroup, ok := r[apiGroup]
								if !ok {
									fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), "API Group", apiGroup, "does not exist")
									mark = red(getEmoji(CROSS_MARK))
									message = fmt.Sprintf(" (API Group '%s' does not exist)", apiGroup)
								} else {
									verbs, ok := availableApiGroup[resourceName]
									if !ok {
										fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), "Resource", resourceName, "does not exist")
										mark = red(getEmoji(CROSS_MARK))
										message = fmt.Sprintf(" (Resource '%s' does not exist)", resourceName)
									} else {
										verbMessage, ok := validateVerbs(rule.Verbs, verbs)
										if !ok {
											mark = red(getEmoji(CROSS_MARK))
											message = verbMessage
										}
									}
								}
								root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)#Role/%s (%s)#%s#%s verbs=%s %s%s",
									sa.Name,
									sa.Namespace,
									roleBinding.Name,
									roleBinding.Namespace,
									role.Name,
									role.Namespace,
									getApiGroup(apiGroup),
									resourceName,
									rule.Verbs,
									mark,
									message))
							}
						}
					}
				}
			} else {
				clusterRole, err := client.RbacV1().ClusterRoles().Get(ctx, roleBinding.RoleRef.Name, metav1.GetOptions{})
				if err != nil {
					fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), err)
					root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)#ClusterRole/%s %s- %s",
						sa.Name,
						sa.Namespace,
						roleBinding.Name,
						roleBinding.Namespace,
						roleBinding.RoleRef.Name,
						getEmoji(CROSS_MARK),
						red("MISSING!!")))
				} else {
					// lets get the permissions
					for _, rule := range clusterRole.Rules {
						for _, resourceName := range rule.Resources {
							for _, apiGroup := range rule.APIGroups {
								mark := green(getEmoji(CHECK_MARK))
								message := ""
								availableApiGroup, ok := r[apiGroup]
								if !ok {
									fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), "API Group", apiGroup, "does not exist")
									mark = red(getEmoji(CROSS_MARK))
									message = fmt.Sprintf(" (API Group '%s' does not exist)", apiGroup)
								} else {
									verbs, ok := availableApiGroup[resourceName]
									if !ok {
										fmt.Println(red(getEmoji(NO_ENTRY)+"WARNING"), "Resource", resourceName, "does not exist")
										mark = red(getEmoji(CROSS_MARK))
										message = fmt.Sprintf(" (Resource '%s' does not exist)", resourceName)
									} else {
										verbMessage, ok := validateVerbs(rule.Verbs, verbs)
										if !ok {
											mark = red(getEmoji(CROSS_MARK))
											message = verbMessage
										}
									}
								}
								root.Add(fmt.Sprintf("ServiceAccount/%s (%s)#RoleBinding/%s (%s)#ClusterRole/%s#%s#%s verbs=%s %s%s",
									sa.Name,
									sa.Namespace,
									roleBinding.Name,
									roleBinding.Namespace,
									clusterRole.Name,
									getApiGroup(apiGroup),
									resourceName,
									rule.Verbs,
									mark,
									message))
							}
						}
					}
				}
			}
		}
	}

	root.Fprint(os.Stdout, true, "")

	return nil
}

func validateVerbs(configuredVerbs metav1.Verbs, availableVerbs []string) (string, bool) {
	var invalid []string
	for _, configuredVerb := range configuredVerbs {
		if !contains(configuredVerb, availableVerbs) {
			invalid = append(invalid, configuredVerb)
		}
	}
	return " (Permissions '" + strings.Join(invalid, ", ") + "' are missing)", len(invalid) == 0
}

func contains(check string, list []string) bool {
	for _, in := range list {
		if in == check {
			return true
		}
	}
	return false
}

func matches(subjects []v1.Subject, namespace string, name string) bool {
	for _, sub := range subjects {
		if sub.Namespace != "" {
			if sub.Kind == "ServiceAccount" && sub.Name == name && sub.Namespace == namespace {
				return true
			}
		} else {
			if sub.Kind == "ServiceAccount" && sub.Name == name {
				return true
			}
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

func green(in string) string {
	if noColor {
		return in
	}
	greenOutput := ansi.ColorFunc("green")
	return greenOutput(in)
}

func red(in string) string {
	if noColor {
		return in
	}
	redOutput := ansi.ColorFunc("red")
	return redOutput(in)
}

func getEmoji(in string) string {
	if noColor {
		switch in {
		case CROSS_MARK:
			return "X "
		case NO_ENTRY:
			return "!! "
		default:
			return emoji.Sprint(in)
		}
	}
	return emoji.Sprint(in)
}
