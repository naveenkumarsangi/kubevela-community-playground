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

// ManagedRecord defines readiness for a resource whose status arrives
// asynchronously. Missing status or conditions evaluate to not ready, while
// Ready and Synced must both be true before the component becomes healthy.
func ManagedRecord() *defkit.ComponentDefinition {
	recordName := defkit.String("recordName").Description("Name of the record to manage")
	zone := defkit.String("zone").Description("Zone the record belongs to")

	h := defkit.Health()

	return defkit.NewComponent("managed-record").
		Description("Manages an external record with status-safe readiness checks").
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
