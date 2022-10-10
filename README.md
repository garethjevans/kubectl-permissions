# permissions

TODO make this a kubectl plugin

## Example

```
❯ ./permissions default
ServiceAccount/default (gareth-dev)
└ ClusterRoleBinding/read-k8s-tint
  └ ClusterRole/k8s-reader-tint
    ├ conventions.carto.run
    │ └ podintents verbs=[get watch list]
    ├
    │ ├ pods verbs=[get watch list]
    │ ├ pods/log verbs=[get watch list]
    │ ├ services verbs=[get watch list]
    │ └ configmaps verbs=[get watch list]
    ├ networking.internal.knative.dev
    │ └ serverlessservices verbs=[get watch list]
    ├ serving.knative.dev
    │ ├ revisions verbs=[get watch list]
    │ ├ routes verbs=[get watch list]
    │ ├ services verbs=[get watch list]
    │ └ configurations verbs=[get watch list]
    ├ carto.run
    │ ├ runnables verbs=[get watch list]
    │ ├ workloads verbs=[get watch list]
    │ ├ clusterconfigtemplates verbs=[get watch list]
    │ ├ clusterdeliveries verbs=[get watch list]
    │ ├ clustersourcetemplates verbs=[get watch list]
    │ ├ clustersupplychains verbs=[get watch list]
    │ ├ clustertemplates verbs=[get watch list]
    │ ├ deliverables verbs=[get watch list]
    │ ├ clusterdeploymenttemplates verbs=[get watch list]
    │ ├ clusterimagetemplates verbs=[get watch list]
    │ └ clusterruntemplates verbs=[get watch list]
    ├ source.toolkit.fluxcd.io
    │ └ gitrepositories verbs=[get watch list]
    ├ apps
    │ ├ deployments verbs=[get watch list]
    │ └ replicasets verbs=[get watch list]
    ├ autoscaling
    │ └ horizontalpodautoscalers verbs=[get watch list]
    ├ source.apps.tanzu.vmware.com
    │ ├ imagerepositories verbs=[get watch list]
    │ └ mavenartifacts verbs=[get watch list]
    ├ networking.k8s.io
    │ └ ingresses verbs=[get watch list]
    ├ scanning.apps.tanzu.vmware.com
    │ ├ sourcescans verbs=[get watch list]
    │ ├ imagescans verbs=[get watch list]
    │ └ scanpolicies verbs=[get watch list]
    ├ autoscaling.internal.knative.dev
    │ └ podautoscalers verbs=[get watch list]
    ├ kpack.io
    │ ├ images verbs=[get watch list]
    │ └ builds verbs=[get watch list]
    ├ tekton.dev
    │ ├ taskruns verbs=[get watch list]
    │ └ pipelineruns verbs=[get watch list]
    └ kappctrl.k14s.io
      └ apps verbs=[get watch list]
```
