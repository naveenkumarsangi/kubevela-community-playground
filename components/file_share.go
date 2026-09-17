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
	defkit.Register(FileShare())
}

// FileShare accepts mount points as a dynamic-key map with typed object values,
// then renders them as a list for the workload API.
//
// OfObject defines the value schema. ForEachMapWithGuarded preserves each key
// and value while building the list and skips the transform when the optional
// source is absent.
func FileShare() *defkit.ComponentDefinition {
	mountPoints := defkit.Map("mountPoints").
		Optional().
		Description("Mount points keyed by name").
		OfObject(
			defkit.String("path").Required(),
			defkit.String("permissions").Optional(),
		)

	return defkit.NewComponent("file-share").
		Description("Renders named mount-point configuration as a workload list").
		Workload("example.com/v1alpha1", "FileShare").
		Params(mountPoints).
		Template(func(tpl *defkit.Template) {
			vela := defkit.VelaCtx()

			// The key becomes a field on the item, which is the whole point: a
			// list of positional entries that still carry their name.
			mounts := defkit.NewArray().ForEachMapWithGuarded(
				mountPoints.IsSet(),
				mountPoints,
				func(entry *defkit.MapEntryBuilder) {
					entry.Set("name", entry.Key())
					entry.Set("path", entry.Value().Field("path"))

					// Per-entry conditionals work the same as on ItemBuilder, so
					// an absent optional field still yields a concrete value.
					entry.IfSet("permissions", func() {
						entry.Set("permissions", entry.Value().Field("permissions"))
					})
					entry.IfNotSet("permissions", func() {
						entry.Set("permissions", defkit.Lit("0755"))
					})
				},
			)

			tpl.Output(defkit.NewResource("example.com/v1alpha1", "FileShare").
				Set("metadata.name", vela.Name()).
				Set("metadata.namespace", vela.Namespace()).
				Set("spec.mountPoints", mounts))
		})
}
