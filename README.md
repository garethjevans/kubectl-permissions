# permissions

A kubectl plugin to display permissions from a service account

## Installation

### Install a development version

```
make build
make install
```

### Install via krew

```
kubectl krew update
kubectl krew install permissions
```

### Manual Installation

```commandline
curl -LO https://github.com/garethjevans/kubectl-permissions/releases/download/v0.0.4/kubectl-permissions_v0.0.4_darwin_amd64.tar.gz && \
    tar -zxvf kubectl-permissions_v0.0.4_darwin_amd64.tar.gz && \
    sudo mv kubectl-permissions /usr/local/bin
```

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

To display the current version of the plugin you can use:

```commandline
❯ kubectl permissions --version
0.0.4
```

