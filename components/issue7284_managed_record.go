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
	defkit.Register(ManagedRecord())
}

// ManagedRecord demonstrates health conditions that survive an absent status
// (kubevela/kubevela#7284).
//
// A freshly applied resource has no status for the first few reconciles. The
// condition comprehensions used to be emitted at the top level and dereferenced
// context.output.status.conditions unconditionally, so CUE propagated bottom
// through the && chain instead of short-circuiting: the policy errored rather
// than reporting "not healthy yet", and every reconcile logged a health-check
// failure for the whole creation window.
//
// The same expression now guards its own preamble. Missing status reads as
// unhealthy, and a condition entry with no type is skipped instead of breaking
// evaluation. Against the five fixtures from the issue -- {}, {status: {}},
// both-true, one-false, and a typeless entry alongside healthy ones -- isHealth
// evaluates to false, false, true, false, true, with no fixture producing
// bottom.
func ManagedRecord() *defkit.ComponentDefinition {
	recordName := defkit.String("recordName").Description("Name of the record to manage")
	zone := defkit.String("zone").Description("Zone the record belongs to")

	h := defkit.Health()

	return defkit.NewComponent("managed-record").
		Description("Manages an external record; demonstrates null-safe health conditions (issue #7284)").
		Workload("example.com/v1alpha1", "Record").
		Params(recordName, zone).
		HealthPolicyExpr(h.And(
			h.Exists("status"),
			h.Exists("status.conditions"),
			h.AllTrue("Ready", "Synced"),
		)).
		Template(func(tpl *defkit.Template) {
			vela := defkit.VelaCtx()

			tpl.Output(defkit.NewResource("example.com/v1alpha1", "Record").
				Set("metadata.name", vela.Name()).
				Set("metadata.namespace", vela.Namespace()).
				Set("spec.recordName", recordName).
				Set("spec.zone", zone))
		})
}
