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
	defkit.Register(ZoneReplica())
}

// ZoneReplica validates a list of replica placements. ValidateValue builds a
// message from the local zone field so a rejected configuration identifies the
// exact value that conflicts with the primary zone.
//
// A fixed message remains appropriate for the separate empty-zone rule.
func ZoneReplica() *defkit.ComponentDefinition {
	primaryZone := defkit.String("primaryZone").Description("Zone that serves writes")

	replicas := defkit.List("replicas").
		Description("Read replicas, one entry per zone").
		WithFields(
			defkit.String("zone").Description("Zone to place the replica in"),
			defkit.Int("readUnits").Optional().Description("Read capacity for this replica"),
		).
		Validators(
			// Dynamic message: the failing zone is interpolated into the label.
			defkit.ValidateValue(defkit.Interpolation(
				defkit.Lit("zone '"),
				defkit.LocalField("zone"),
				defkit.Lit("' must not match the primary zone"),
			)).
				FailWhen(defkit.Eq(defkit.LocalField("zone"), defkit.Reference("parameter.primaryZone"))).
				WithName("_validateReplicaZone"),

			// Fixed message, for contrast: the rule needs no value to be clear.
			defkit.Validate("replica zone must not be empty").
				FailWhen(defkit.LocalField("zone").IsEmpty()).
				WithName("_validateZoneNotEmpty"),
		)

	return defkit.NewComponent("zone-replica").
		Description("Places read replicas across zones with value-specific validation").
		Workload("example.com/v1alpha1", "ReplicaSet").
		Params(primaryZone, replicas).
		Template(func(tpl *defkit.Template) {
			vela := defkit.VelaCtx()

			items := defkit.NewArray().ForEachWith(replicas, func(item *defkit.ItemBuilder) {
				item.Set("zone", item.Var().Field("zone"))
				item.IfSet("readUnits", func() {
					item.Set("capacity.readUnits", item.Var().Field("readUnits"))
				})
			})

			tpl.Output(defkit.NewResource("example.com/v1alpha1", "ReplicaSet").
				Set("metadata.name", vela.Name()).
				Set("metadata.namespace", vela.Namespace()).
				Set("spec.primaryZone", primaryZone).
				Set("spec.replicas", items))
		})
}
