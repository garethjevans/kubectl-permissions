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

The plugin also has the ability to display any secrets attached to a service account, either as a `secret` or an `imagePullSecret`.

```commandLine
❯ kubectl permissions my-sa --include-secrets
ImagePullSecrets
└ registries-credentials ✔
  └ type=kubernetes.io/dockerconfigjson
Secrets
├ git-ssh ✔
│ ├ tekton.dev/git-0=https://my-git-server
│ └ type=kubernetes.io/basic-auth
└ registry-credentials ✔
  ├ tekton.dev/docker-0=docker.io/my-docker-registry
  └ type=kubernetes.io/dockerconfigjson
```

To display the current version of the plugin you can use:

```commandline
❯ kubectl permissions --version
0.0.4
```

## Verifying the artifacts

All artifacts are checksummed, and the `checksum.txt` file is signed using [cosign](https://github.com/sigstore/cosign).

1. Download your required binary, and also the certificate(`checksums.txt.pem`), signature(`checksums.txt.sig`) and the `checksums.txt` (_this file contains a checksum for each artifact_) file.

```sh
VERSION=v0.0.6
https://github.com/garethjevans/kubectl-permissions/releases/download/$VERSION/checksums.txt.pem
https://github.com/garethjevans/kubectl-permissions/releases/download/$VERSION/checksums.txt.sig
https://github.com/garethjevans/kubectl-permissions/releases/download/$VERSION/checksums.txt
```

2. Now you can verify the signature:

```sh
cosign verify-blob \
  --cert checksums.txt.pem \
  --signature checksums.txt.sig \
  checksums.txt
```
3. To wrap up, you can verify the SHA256 checksums match the downloaded binary:

```sh
sha256sum --ignore-missing -c checksums.txt
```
