package roles

import (
	"k8s.io/client-go/discovery"
	"strings"
)

func DiscoverRolesAndPermissions(d *discovery.DiscoveryClient) (map[string]map[string][]string, error) {
	rolesAndPermissions := make(map[string]map[string][]string)
	_, resourceListArr, err := d.ServerGroupsAndResources()
	if err != nil {
		return rolesAndPermissions, err
	}

	for _, resourceList := range resourceListArr {
		// rbac rules only look at API group names, not name & version
		groupOnly := strings.Split(resourceList.GroupVersion, "/")[0]
		// core API doesn't have a group "name". We set to "core" and replace at the end with a blank string in the rbac policy rule
		if resourceList.GroupVersion == "v1" {
			groupOnly = ""
		}

		_, ok := rolesAndPermissions[groupOnly]
		if !ok {
			rolesAndPermissions[groupOnly] = make(map[string][]string)
		}

		for _, resource := range resourceList.APIResources {
			verbs := make([]string, 0)
			for _, v := range resource.Verbs {
				verbs = append(verbs, v)
			}

			// always add '*' as a wildcard
			verbs = append(verbs, "*")

			rolesAndPermissions[groupOnly][resource.Name] = verbs
		}
	}

	return rolesAndPermissions, nil
}
