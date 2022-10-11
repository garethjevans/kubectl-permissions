package integration

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"testing"
	"time"
)

func TestPluginIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	gomega.RegisterTestingT(t)

	command := exec.Command("kubectl", "permissions", "sa-under-test", "-n", "test-namespace")
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	if err != nil {
		t.Errorf("unable to exec command %s", err)
	}

	response := string(session.Wait(10 * time.Second).Out.Contents())

	expected := `[WARNING] roles.rbac.authorization.k8s.io "a-missing-role" not found
ServiceAccount/sa-under-test (test-namespace)
├ ClusterRoleBinding/cluster-roles
│ └ ClusterRole/cluster-level-role
│   ├ <default>
│   │ ├ configmaps verbs=[get watch list]
│   │ ├ pods verbs=[get watch list]
│   │ ├ pods/log verbs=[get watch list]
│   │ └ services verbs=[get watch list]
│   ├ apps
│   │ ├ deployments verbs=[get watch list]
│   │ └ replicasets verbs=[get watch list]
│   └ networking.k8s.io
│     └ ingresses verbs=[get]
├ RoleBinding/missconfigured (test-namespace)
│ └ Role/a-missing-role (a-missing-role) MISSING!!
└ RoleBinding/namespaced-roles (test-namespace)
  └ Role/namespaced-role (test-namespace)
    ├ kpack.io
    │ ├ builds verbs=[get watch list]
    │ └ images verbs=[get watch list]
    ├ source.toolkit.fluxcd.io
    │ └ gitrepositories verbs=[get watch list]
    └ tekton.dev
      ├ pipelineruns verbs=[get watch list]
      └ taskruns verbs=[get watch list]`

	gomega.Expect(response).To(gomega.Equal(expected))
}
