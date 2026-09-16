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

// FileShare demonstrates guarded map-to-list comprehensions
// (kubevela/kubevela#7288).
//
// The input is a name-keyed map and the template needs a positional list that
// keeps the name. ForEachWithVar exposes one iteration variable, so it can reach
// the value but never the key; ForEachMap has both variables but always emits a
// struct. ForEachMapWithGuarded closes the gap: both variables, a list result,
// and a guard for an unset source.
//
// Generated shape:
//
//	[ if parameter["mountPoints"] != _|_ for k, v in parameter.mountPoints { ... } ]
func FileShare() *defkit.ComponentDefinition {
	mountPoints := defkit.Map("mountPoints").
		Optional().
		Description("Mount points keyed by name").
		WithSchema(`[string]: {
	path: string
	permissions?: string
}`)

	return defkit.NewComponent("file-share").
		Description("Exposes named mount points; demonstrates map-to-list comprehensions (issue #7288)").
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
