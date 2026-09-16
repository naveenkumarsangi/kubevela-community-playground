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
	defkit.Register(QueueSet())
}

// QueueSet demonstrates nested field paths inside a per-item ItemBuilder
// (kubevela/kubevela#7282).
//
// ItemBuilder.Set now expands a dot-separated path into nested CUE structs, the
// same way Resource.Set does. Before the fix the whole string became one field
// label, so "throughput.maxReadOps" generated `throughput.maxReadOps: ...`,
// which CUE rejects with "missing ',' in struct literal" — and only at cue vet
// time, never during generation.
//
// The example exercises the three shapes the issue asked to cover: a plain
// nested Set, nested Sets inside If/IfSet/IfNotSet, and sibling writes under one
// parent that have to unify into a single struct.
func QueueSet() *defkit.ComponentDefinition {
	queues := defkit.List("queues").
		Description("Queues to provision, one entry per queue").
		WithFields(
			defkit.String("name").Description("Queue name"),
			defkit.Int("maxReadOps").Optional().Description("Read operations allowed per second"),
			defkit.Int("maxWriteOps").Optional().Description("Write operations allowed per second"),
			defkit.Int("retentionHours").Optional().Description("How long messages are kept, in hours"),
		)

	return defkit.NewComponent("queue-set").
		Description("Provisions a set of queues; demonstrates nested paths in ItemBuilder (issue #7282)").
		Workload("example.com/v1alpha1", "QueueSet").
		Params(queues).
		Template(func(tpl *defkit.Template) {
			vela := defkit.VelaCtx()

			items := defkit.NewArray().ForEachWith(queues, func(item *defkit.ItemBuilder) {
				// Plain nested path: one Set, two levels deep.
				item.Set("spec.queueName", item.Var().Field("name"))

				// Sibling nested writes under a shared parent. Each lands in its
				// own if block, and the two blocks unify into one throughput
				// struct rather than colliding.
				item.IfSet("maxReadOps", func() {
					item.Set("spec.throughput.maxReadOps", item.Var().Field("maxReadOps"))
				})
				item.IfSet("maxWriteOps", func() {
					item.Set("spec.throughput.maxWriteOps", item.Var().Field("maxWriteOps"))
				})

				// Nested path in the negative branch too, so an unset input still
				// produces a concrete field.
				item.IfNotSet("retentionHours", func() {
					item.Set("spec.retention.hours", defkit.Lit(24))
				})
				item.IfSet("retentionHours", func() {
					item.Set("spec.retention.hours", item.Var().Field("retentionHours"))
				})
			})

			tpl.Output(defkit.NewResource("example.com/v1alpha1", "QueueSet").
				Set("metadata.name", vela.Name()).
				Set("metadata.namespace", vela.Namespace()).
				Set("spec.queues", items))
		})
}
