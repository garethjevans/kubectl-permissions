package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func TestPluginIntegration(t *testing.T) {
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

	expected := "\x1b[0;31m⛔ WARNING\x1b[0m roles.rbac.authorization.k8s.io \"a-missing-role\" not found\nServiceAccount/sa-under-test (test-namespace)\n\x1b[0;94;40m├\x1b[0m ClusterRoleBinding/cluster-roles\n\x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m ClusterRole/cluster-level-role\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m├\x1b[0m apps\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m deployments verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m replicasets verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m├\x1b[0m core.k8s.io\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m configmaps verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m pods verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m pods/log verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m services verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n\x1b[0;94;40m│\x1b[0m   \x1b[0;94;40m└\x1b[0m networking.k8s.io\n\x1b[0;94;40m│\x1b[0m     \x1b[0;94;40m└\x1b[0m ingresses verbs=[get] \x1b[0;32m✔ \x1b[0m\n\x1b[0;94;40m├\x1b[0m RoleBinding/missconfigured (test-namespace)\n\x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m Role/a-missing-role (a-missing-role) ❌ - \x1b[0;31mMISSING!!\x1b[0m\n\x1b[0;94;40m└\x1b[0m RoleBinding/namespaced-roles (test-namespace)\n  \x1b[0;94;40m└\x1b[0m Role/namespaced-role (test-namespace)\n    \x1b[0;94;40m├\x1b[0m kpack.io\n    \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m├\x1b[0m builds verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n    \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m images verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n    \x1b[0;94;40m├\x1b[0m source.toolkit.fluxcd.io\n    \x1b[0;94;40m│\x1b[0m \x1b[0;94;40m└\x1b[0m gitrepositories verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n    \x1b[0;94;40m└\x1b[0m tekton.dev\n      \x1b[0;94;40m├\x1b[0m pipelineruns verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n      \x1b[0;94;40m└\x1b[0m taskruns verbs=[get watch list] \x1b[0;32m✔ \x1b[0m\n"

	Expect(strings.TrimSpace(response)).To(Equal(strings.TrimSpace(expected)))
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
	expected := `!! WARNING roles.rbac.authorization.k8s.io "a-missing-role" not found
ServiceAccount/sa-under-test (test-namespace)
├ ClusterRoleBinding/cluster-roles
│ └ ClusterRole/cluster-level-role
│   ├ apps
│   │ ├ deployments verbs=[get watch list] ✔ 
│   │ └ replicasets verbs=[get watch list] ✔ 
│   ├ core.k8s.io
│   │ ├ configmaps verbs=[get watch list] ✔ 
│   │ ├ pods verbs=[get watch list] ✔ 
│   │ ├ pods/log verbs=[get watch list] ✔ 
│   │ └ services verbs=[get watch list] ✔ 
│   └ networking.k8s.io
│     └ ingresses verbs=[get] ✔ 
├ RoleBinding/missconfigured (test-namespace)
│ └ Role/a-missing-role (a-missing-role) X - MISSING!!
└ RoleBinding/namespaced-roles (test-namespace)
  └ Role/namespaced-role (test-namespace)
    ├ kpack.io
    │ ├ builds verbs=[get watch list] ✔ 
    │ └ images verbs=[get watch list] ✔ 
    ├ source.toolkit.fluxcd.io
    │ └ gitrepositories verbs=[get watch list] ✔ 
    └ tekton.dev
      ├ pipelineruns verbs=[get watch list] ✔ 
      └ taskruns verbs=[get watch list] ✔ 
`
	Expect(strings.TrimSpace(response)).To(Equal(strings.TrimSpace(expected)))
}
