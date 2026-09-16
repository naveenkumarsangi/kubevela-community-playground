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
	defkit.Register(TenantSpace())
}

// TenantSpace demonstrates negative regex conditions
// (kubevela/kubevela#7353).
//
// Defkit had typed helpers for positive runtime regex matching but nothing
// symmetric for the negative case, even though CUE has a native `!~` operator.
// Callers had to wrap a positive match in Not(...), which is semantically fine
// but less discoverable and emits `!(x =~ "p")` instead of `x !~ "p"`.
//
// NotMatches now exists on both fluent condition builders -- StringParam for
// runtime parameter conditions, LocalFieldRef for validators -- and reads
// naturally in the case the issue called out: a validator that should fail when
// a value does not match an allowed pattern.
func TenantSpace() *defkit.ComponentDefinition {
	tenantName := defkit.String("tenantName").Description("Tenant that owns this space")
	displayName := defkit.String("displayName").Optional().Description("Human-readable name")

	return defkit.NewComponent("tenant-space").
		Description("Allocates a namespaced tenant space; demonstrates negative regex conditions (issue #7353)").
		Workload("example.com/v1alpha1", "TenantSpace").
		Params(tenantName, displayName).
		Validators(
			// The case from the issue: fail when the value does not match the
			// allowed pattern. NotMatches states that directly, where the old
			// form was Not(LocalField("tenantName").Matches(...)).
			defkit.Validate("tenantName contains unsupported characters").
				FailWhen(defkit.LocalField("tenantName").NotMatches(`^[a-z0-9-]+$`)).
				WithName("_validateTenantName"),

			// And the plainly negative rule, which needed no wrapper either way.
			defkit.Validate("tenantName must not end with a hyphen").
				FailWhen(defkit.LocalField("tenantName").Matches(`.*-$`)).
				WithName("_validateTenantSuffix"),
		).
		Template(func(tpl *defkit.Template) {
			vela := defkit.VelaCtx()

			res := defkit.NewResource("example.com/v1alpha1", "TenantSpace").
				Set("metadata.name", vela.Name()).
				Set("metadata.namespace", vela.Namespace()).
				Set("spec.tenantName", tenantName)

			// Runtime condition on a parameter, the StringParam half of the pair.
			// Reserved tenants are the ones outside the customer prefix.
			res.SetIf(tenantName.NotMatches(`^cust-`), "spec.tier", defkit.Lit("internal"))
			res.SetIf(tenantName.Matches(`^cust-`), "spec.tier", defkit.Lit("customer"))
			res.SetIf(displayName.IsSet(), "spec.displayName", displayName)

			tpl.Output(res)
		})
}
