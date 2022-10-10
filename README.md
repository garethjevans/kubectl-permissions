# permissions

A kubectl plugin to display permissions from a service account

## Installation

```
make build
make install
```

TODO make this `krew` installable

## Example

```
❯ kubectl permissions default
ServiceAccount/default (tap-gui)
└ ClusterRoleBinding/read-k8s-tint
  └ ClusterRole/k8s-reader-tint
    ├ <default>
    │ ├ configmaps verbs=[get watch list]
    │ ├ pods verbs=[get watch list]
    │ ├ pods/log verbs=[get watch list]
    │ └ services verbs=[get watch list]
    ├ apps
    │ ├ deployments verbs=[get watch list]
    │ └ replicasets verbs=[get watch list]
    ├ autoscaling
    │ └ horizontalpodautoscalers verbs=[get watch list]
    ├ autoscaling.internal.knative.dev
    │ └ podautoscalers verbs=[get watch list]
    ├ carto.run
    │ ├ clusterconfigtemplates verbs=[get watch list]
    │ ├ clusterdeliveries verbs=[get watch list]
    │ ├ clusterdeploymenttemplates verbs=[get watch list]
    │ ├ clusterimagetemplates verbs=[get watch list]
    │ ├ clusterruntemplates verbs=[get watch list]
    │ ├ clustersourcetemplates verbs=[get watch list]
    │ ├ clustersupplychains verbs=[get watch list]
    │ ├ clustertemplates verbs=[get watch list]
    │ ├ deliverables verbs=[get watch list]
    │ ├ runnables verbs=[get watch list]
    │ └ workloads verbs=[get watch list]
    ├ conventions.carto.run
    │ └ podintents verbs=[get watch list]
    ├ kappctrl.k14s.io
    │ └ apps verbs=[get watch list]
    ├ kpack.io
    │ ├ builds verbs=[get watch list]
    │ └ images verbs=[get watch list]
    ├ networking.internal.knative.dev
    │ └ serverlessservices verbs=[get watch list]
    ├ networking.k8s.io
    │ └ ingresses verbs=[get watch list]
    ├ scanning.apps.tanzu.vmware.com
    │ ├ imagescans verbs=[get watch list]
    │ ├ scanpolicies verbs=[get watch list]
    │ └ sourcescans verbs=[get watch list]
    ├ serving.knative.dev
    │ ├ configurations verbs=[get watch list]
    │ ├ revisions verbs=[get watch list]
    │ ├ routes verbs=[get watch list]
    │ └ services verbs=[get watch list]
    ├ source.apps.tanzu.vmware.com
    │ ├ imagerepositories verbs=[get watch list]
    │ └ mavenartifacts verbs=[get watch list]
    ├ source.toolkit.fluxcd.io
    │ └ gitrepositories verbs=[get watch list]
    └ tekton.dev
      ├ pipelineruns verbs=[get watch list]
      └ taskruns verbs=[get watch list]
```
