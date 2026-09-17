/*
Copyright 2025 The KubeVela Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package components

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

func init() {
	defkit.Register(RouteCatalog())
}

// RouteCatalog demonstrates structured ForEachMap body operations
// (kubevela/kubevela#7283).
//
// Each input entry maps a route name to a backend. WithBody turns that scalar
// value into a structured object under the original key, and the generator now
// renders those stored body operations for every map entry.
func RouteCatalog() *defkit.ComponentDefinition {
	routes := defkit.StringKeyMap("routes").
		Description("Backend service keyed by route name")

	return defkit.NewComponent("route-catalog").
		Description("Builds structured route entries from a string map; demonstrates ForEachMap body rendering (issue #7283)").
		Workload("example.com/v1alpha1", "RouteCatalog").
		Params(routes).
		Template(func(tpl *defkit.Template) {
			body := defkit.NewResource("", "").
				Set("service", defkit.Reference("backendValue")).
				Set("metadata.sourceRoute", defkit.Reference("routeKey")).
				Ops()
			catalog := defkit.ForEachMap().
				Over("parameter.routes").
				WithVars("routeKey", "backendValue").
				WithBody(body...)

			tpl.Output(defkit.NewResource("example.com/v1alpha1", "RouteCatalog").
				Set("spec.routes", catalog))
		})
}
