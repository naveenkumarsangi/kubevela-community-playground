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

// Package main renders every registered definition to CUE on disk, so the
// generated text can be inspected and checked with `cue vet`.
//
// Usage: go run ./cmd/generate [output-dir]   (default: generated)
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"

	// Import packages to trigger init() registration
	_ "github.com/myorg/playground/components"
	_ "github.com/myorg/playground/policies"
	_ "github.com/myorg/playground/traits"
	_ "github.com/myorg/playground/workflowsteps"
)

var subdirs = map[defkit.DefinitionType]string{
	defkit.DefinitionTypeComponent:    "component",
	defkit.DefinitionTypeTrait:        "trait",
	defkit.DefinitionTypePolicy:       "policy",
	defkit.DefinitionTypeWorkflowStep: "workflowstep",
}

func main() {
	outputDir := "generated"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	defs := defkit.All()
	if len(defs) == 0 {
		fmt.Fprintln(os.Stderr, "no definitions registered")
		os.Exit(1)
	}

	for _, def := range defs {
		subdir, ok := subdirs[def.DefType()]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown definition type %q for %q, skipping\n", def.DefType(), def.DefName())
			continue
		}
		dir := filepath.Join(outputDir, subdir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create %s: %v\n", dir, err)
			os.Exit(1)
		}
		path := filepath.Join(dir, def.DefName()+".cue")
		if err := os.WriteFile(path, []byte(def.ToCue()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Println(path)
	}
}
