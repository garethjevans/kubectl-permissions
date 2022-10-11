# permissions

A kubectl plugin to display permissions from a service account

## Installation

```
make build
make install
```

TODO make this `krew` installable

## Example

Based on the roles configured in the example-rbac.yaml:

```commandLine
❯ kubectl permissions sa-under-test -n test-namespace
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
```

The plugin will also highlight when configured roles are missing:

```commandLine
❯ kubectl permissions sa-under-test -n test-namespace
⛔  WARNING roles.rbac.authorization.k8s.io "a-missing-role" not found
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
│ └ Role/a-missing-role (a-missing-role) ❌ - MISSING!!
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
```
