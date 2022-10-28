# kubectl-permissions

![GitHub all releases](https://img.shields.io/github/downloads/garethjevans/kubectl-permissions/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/garethjevans/kubectl-permissions)
[![Go Report Card](https://goreportcard.com/badge/github.com/garethjevans/kubectl-permissions)](https://goreportcard.com/report/github.com/garethjevans/kubectl-permissions)

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
│   │ ├ pods/log verbs=[get] ✔
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
❯ kubectl permissions invalid-sa
⛔ WARNING roles.rbac.authorization.k8s.io "missing-role" not found
⛔ WARNING API Group bingbong.io does not exist
⛔ WARNING Resource invalid does not exist
ServiceAccount/invalid-sa (test-namespace)
├ RoleBinding/missing-role-binding (test-namespace)
│ └ Role/missing-role (missing-role) ❌ - MISSING!!
└ RoleBinding/missing-role-binding2 (test-namespace)
  └ Role/invalid-role (test-namespace)
    ├ bingbong.io
    │ └ something verbs=[get watch list] ❌  (API Group 'bingbong.io' does not exist)
    ├ source.toolkit.fluxcd.io
    │ └ gitrepositories verbs=[laugh] ❌  (Permissions 'laugh' are missing)
    └ tekton.dev
      └ invalid verbs=[get] ❌  (Resource 'invalid' does not exist)
```

To display the current version of the plugin you can use:

```commandline
❯ kubectl permissions --version
0.0.4
```

