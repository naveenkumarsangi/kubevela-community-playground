# DefKit: from component work to upstream improvements

This repository supports a KubeVela community-call demo. It contains small,
runnable examples of DefKit capabilities added while we were building new OAM
components and migrating selected hand-written CUE definitions to Go.

We began with a practical goal: use DefKit for real component delivery. As we
applied it to more component shapes, we identified opportunities to broaden the
typed authoring experience. We reduced those cases to focused examples and
worked with the community on reusable upstream enhancements.

## Why we chose DefKit

DefKit is KubeVela's Go SDK for authoring X-Definitions. The author writes typed,
fluent Go; DefKit generates the CUE definition that KubeVela already understands.
It gave us several practical advantages:

- the same Go tooling, review workflow, and IDE support used for controllers;
- reusable builders instead of repeating large CUE fragments;
- ordinary Go tests and module versioning around definition authoring;
- a clearer migration path for teams that maintain Go and CUE together; and
- standard CUE output, so adopting DefKit does not change the KubeVela runtime.

CUE remains the generated language and KubeVela's evaluation engine. DefKit adds
a maintainable Go authoring layer while preserving raw CUE as an option for
specialized cases.

## What the component work helped improve

We used DefKit both to create new OAM components and to migrate existing
components that had been maintained directly in CUE. That experience highlighted
three areas where DefKit could support a wider range of component designs:

1. **Authoring expressiveness** — typed builders now cover more nested paths,
   map-to-list transformations, and negative regex conditions.
2. **Feedback and correctness** — validators can identify the value being
   rejected, while health logic treats a newly created resource with no status
   as a normal not-ready state.
3. **Multi-resource lifecycle** — health builders can describe a primary
   resource together with named and repeated auxiliary outputs.

We initially opened eleven focused DefKit issues with reproductions and expected
behaviour. A later migration case added the negative-regex enhancement. Some
changes were contributed by us; others were implemented by community
contributors.

## Merged outcomes demonstrated here

As of 17 September 2026, these six improvements are merged:

| Area | Issue and merged PR | What changed | Example |
| --- | --- | --- | --- |
| Nested structures | [#7282](https://github.com/kubevela/kubevela/issues/7282) / [#7302](https://github.com/kubevela/kubevela/pull/7302) | `ItemBuilder.Set` expands dotted paths instead of emitting an invalid field label | `queue-set` |
| Early reconciliation | [#7284](https://github.com/kubevela/kubevela/issues/7284) / [#7299](https://github.com/kubevela/kubevela/pull/7299) | Missing status evaluates as not healthy instead of producing CUE bottom | `managed-record` |
| Collection transforms | [#7288](https://github.com/kubevela/kubevela/issues/7288) / [#7330](https://github.com/kubevela/kubevela/pull/7330) | Map entries can become list items while retaining both key and value | `file-share` |
| Actionable validation | [#7289](https://github.com/kubevela/kubevela/issues/7289) / [#7348](https://github.com/kubevela/kubevela/pull/7348) | Validator messages can interpolate the value that failed | `zone-replica` |
| Multi-output health | [#7290](https://github.com/kubevela/kubevela/issues/7290) / [#7332](https://github.com/kubevela/kubevela/pull/7332) | Health expressions can target named outputs and aggregate output groups | `composite-store` |
| Fluent conditions | [#7353](https://github.com/kubevela/kubevela/issues/7353) / [#7352](https://github.com/kubevela/kubevela/pull/7352) | `NotMatches` emits CUE's native negative-regex operator | `tenant-space` |

The repository keeps one deliberately small component per capability. They are
neutral examples rather than copies of production definitions, so the generated
CUE and dry-run output remain concise enough for a short demo.

## Quick start

Requirements:

- Go 1.23.8 or later;
- a `vela` CLI with `def vet` and offline `dry-run` support; and
- network access during the first `go mod download`.

The Go module pins a post-merge KubeVela revision that contains all six APIs.
No Kubernetes cluster is needed.

```bash
go mod download
./scripts/verify-demo.sh
```

The verification script:

1. generates all six ComponentDefinitions;
2. compiles every Go package;
3. validates every generated CUE definition;
4. dry-runs every valid Application; and
5. confirms that the intentionally invalid replica example fails with the
   offending zone in its message.

To run the presentation steps manually:

```bash
go run ./cmd/generate generated

vela dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml

vela dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml

vela dry-run --offline \
  -d generated/component \
  -f examples/composite-store.yaml
```

The second command is expected to fail. For this part of the demo, focus on this
output:

```text
parameter.replicas.1._validateReplicaZone."zone 'us-east-1' must not match the primary zone": conflicting values true and false
```

See [`PRESENTATION.md`](PRESENTATION.md) for the talk track and
[`DEMO-RUNBOOK.md`](DEMO-RUNBOOK.md) for the exact live sequence, handoff cues,
and recovery commands.

## Repository layout

```text
components/            Go-authored DefKit ComponentDefinitions
examples/              Applications used by offline dry-run
cmd/generate/           writes registered definitions as CUE
cmd/register/           emits the DefKit registry as JSON
scripts/verify-demo.sh  repeatable pre-call verification
PRESENTATION.md         timed narrative for the community call
DEMO-RUNBOOK.md         live commands and expected signals
```

## Scope of the demo

The examples emit resources under `example.com/v1alpha1`. They intentionally do
not require CRDs or controllers because the demo is about the definition-authoring
path: Go source, generated CUE, rendered resources, validation feedback, and
health-policy structure. A live cluster would add setup and failure modes without
making those DefKit behaviours clearer.
