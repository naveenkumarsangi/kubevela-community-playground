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
	defkit.Register(CompositeStore())
}

// CompositeStore emits a primary store, a named access policy, and a group of
// access-point outputs.
//
// At scopes a health expression to a specific output. Every applies the same
// expression to all outputs with a shared prefix. The component becomes healthy
// only when the store, policy, and every requested access point report Ready.
func CompositeStore() *defkit.ComponentDefinition {
	storeName := defkit.String("storeName").Description("Name of the store to provision")
	accessPoints := defkit.Int("accessPoints").
		Default(0).
		Description("How many access points to create, 0, 1 or 2")

	h := defkit.Health()
	vela := defkit.VelaCtx()

	return defkit.NewComponent("composite-store").
		Description("Provisions a store and evaluates health across all related resources").
		Workload("example.com/v1alpha1", "Store").
		Params(storeName, accessPoints).
		HealthPolicyExpr(h.And(
			// Default scope: the primary output.
			h.At(vela.Output()).Condition("Ready").IsTrue(),
			// A named auxiliary output.
			h.At(vela.Outputs("accessPolicy")).Condition("Ready").IsTrue(),
			// Every output sharing the accessPoint prefix. No AllowEmpty(), so a
			// component with nothing created yet is unhealthy.
			h.Every(defkit.OutputsWithPrefix("accessPoint"), func(item *defkit.HealthScope) defkit.HealthExpression {
				return item.Condition("Ready").IsTrue()
			}),
		)).
		Template(func(tpl *defkit.Template) {
			tpl.Output(defkit.NewResource("example.com/v1alpha1", "Store").
				Set("metadata.name", vela.Name()).
				Set("metadata.namespace", vela.Namespace()).
				Set("spec.storeName", storeName))

			tpl.Outputs("accessPolicy", defkit.NewResource("example.com/v1alpha1", "AccessPolicy").
				Set("metadata.name", defkit.Interpolation(vela.Name(), defkit.Lit("-policy"))).
				Set("metadata.namespace", vela.Namespace()).
				Set("spec.storeRef", vela.Name()))

			// Numbered outputs sharing one prefix. Output names are always
			// written out -- there is no API that generates them from a count --
			// so the slots are enumerated and each is guarded by the request.
			tpl.OutputsIf(defkit.Ge(accessPoints, defkit.Lit(1)), "accessPoint1",
				accessPointResource(vela, "-ap-1"))
			tpl.OutputsIf(defkit.Ge(accessPoints, defkit.Lit(2)), "accessPoint2",
				accessPointResource(vela, "-ap-2"))
		})
}

func accessPointResource(vela *defkit.VelaContext, suffix string) *defkit.Resource {
	return defkit.NewResource("example.com/v1alpha1", "AccessPoint").
		Set("metadata.name", defkit.Interpolation(vela.Name(), defkit.Lit(suffix))).
		Set("metadata.namespace", vela.Namespace()).
		Set("spec.storeRef", vela.Name())
}
