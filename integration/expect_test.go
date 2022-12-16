package integration

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type test struct {
		name  string
		sa    string
		flags []string

		expectedOut string
		expectedErr string
	}

	tests := []test{
		{
			name:  "sa-under-test",
			sa:    "sa-under-test",
			flags: []string{"-n", "test-namespace"},

			expectedOut: "ServiceAccount/sa-under-test (test-namespace)\n" +
				"\x1b[0;94;40m├\x1b[0m ClusterRoleBinding/cluster-roles\n" +
				"\x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m ClusterRole/cluster-level-role\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m├\x1b[0m apps\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m deployments verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m replicasets verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m├\x1b[0m core.k8s.io\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m configmaps verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m pods verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m pods/log verbs=[get] \x1b[0;32m✔ \x1b[0m\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m services verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n" +
				"\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m└\x1b[0m networking.k8s.io\n\x1b[0;94;40m│\x1b[0m     \x1b[0;94;40m└\x1b[0m ingresses verbs=[get] \x1b[0;32m✔ \x1b[0m\n" +
				"\x1b[0;94;40m└\x1b[0m RoleBinding/namespaced-roles (test-namespace)\n" +
				"  \x1b[0;94;40m└\x1b[0m Role/namespaced-role (test-namespace)\n" +
				"    \x1b[0;94;40m└\x1b[0m core.k8s.io\n" +
				"      \x1b[0;94;40m└\x1b[0m secrets verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n",
		},

		{
			name:  "monitoring",
			sa:    "monitoring",
			flags: []string{"-n", "test-namespace"},

			expectedOut: "ServiceAccount/monitoring (test-namespace)\n" +
				"\x1b[0;94;40m└\x1b[0m ClusterRoleBinding/monitoring\n" +
				"  \x1b[0;94;40m└\x1b[0m ClusterRole/monitoring\n" +
				"    \x1b[0;94;40m└\x1b[0m core.k8s.io\n" +
				"      \x1b[0;94;40m├\x1b[0m endpoints verbs=[create] \x1b[0;32m✔ \x1b[0m\n" +
				"      \x1b[0;94;40m├\x1b[0m endpoints verbs=[get list watch] \x1b[0;32m✔ \x1b[0m\n" +
				"      \x1b[0;94;40m├\x1b[0m pods verbs=[create] \x1b[0;32m✔ \x1b[0m\n" +
				"      \x1b[0;94;40m├\x1b[0m pods verbs=[get list watch] \x1b[0;32m✔ \x1b[0m\n" +
				"      \x1b[0;94;40m├\x1b[0m services verbs=[create] \x1b[0;32m✔ \x1b[0m\n" +
				"      \x1b[0;94;40m└\x1b[0m services verbs=[get list watch] \x1b[0;32m✔ \x1b[0m\n",
		},

		{
			name:  "sa-that-doesnt-exist",
			sa:    "sa-that-doesnt-exist",
			flags: []string{"-n", "test-namespace"},

			expectedOut: "",
			expectedErr: "Error: serviceaccounts \"sa-that-doesnt-exist\" not found",
		},

		{
			name:  "sa-with-no-bindings",
			sa:    "sa-with-no-bindings",
			flags: []string{"-n", "test-namespace"},

			expectedOut: "ServiceAccount/sa-with-no-bindings (test-namespace)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterTestingT(t)

			var args []string
			args = append(args, "permissions", tt.sa)
			args = append(args, tt.flags...)

			command := exec.Command("kubectl", args...)
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			if err != nil {
				t.Errorf("unable to exec command %s", err)
			}

			session = session.Wait(10 * time.Second)
			stdOut := string(session.Out.Contents())
			stdErr := string(session.Err.Contents())

			assert.Equal(t, strings.TrimSpace(tt.expectedOut), strings.TrimSpace(stdOut))
			assert.Equal(t, strings.TrimSpace(tt.expectedErr), strings.TrimSpace(stdErr))
		})
	}
}

func TestPluginIntegrationNoColor(t *testing.T) {
	os.Setenv("NO_COLOR", "true")
	defer os.Unsetenv("NO_COLOR")
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	RegisterTestingT(t)

	command := exec.Command("kubectl", "permissions", "sa-under-test", "-n", "test-namespace")
	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	if err != nil {
		t.Errorf("unable to exec command %s", err)
	}

	response := string(session.Wait(10 * time.Second).Out.Contents())
	expected := `ServiceAccount/sa-under-test (test-namespace)
├ ClusterRoleBinding/cluster-roles
│ └ ClusterRole/cluster-level-role
│   ├ apps
│   │ ├ deployments verbs=[get watch list] ✔ 
│   │ └ replicasets verbs=[get watch list] ✔ 
│   ├ core.k8s.io
│   │ ├ configmaps verbs=[get watch list] ✔ 
│   │ ├ pods verbs=[get watch list] ✔ 
│   │ ├ pods/log verbs=[get] ✔ 
│   │ └ services verbs=[get watch list] ✔ 
│   └ networking.k8s.io
│     └ ingresses verbs=[get] ✔ 
└ RoleBinding/namespaced-roles (test-namespace)
  └ Role/namespaced-role (test-namespace)
    └ core.k8s.io
      └ secrets verbs=[get watch list] ✔ 
`
	Expect(strings.TrimSpace(response)).To(Equal(strings.TrimSpace(expected)))
}
