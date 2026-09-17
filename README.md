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

1. **Authoring expressiveness** — typed builders now cover structured dynamic
   maps, nested paths, map-to-map and map-to-list transformations, and negative
   regex conditions.
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

As of 17 September 2026, eight of the twelve improvements are merged:

| Area | Issue and merged PR | What changed | Example |
| --- | --- | --- | --- |
| Nested structures | [#7282](https://github.com/kubevela/kubevela/issues/7282) / [#7302](https://github.com/kubevela/kubevela/pull/7302) | `ItemBuilder.Set` expands dotted paths into nested fields | `queue-set` |
| Structured map transforms | [#7283](https://github.com/kubevela/kubevela/issues/7283) / [#7339](https://github.com/kubevela/kubevela/pull/7339) | `ForEachMap.WithBody` renders structured values for each map entry | `route-catalog` |
| Early reconciliation | [#7284](https://github.com/kubevela/kubevela/issues/7284) / [#7299](https://github.com/kubevela/kubevela/pull/7299) | Missing status evaluates as not healthy instead of producing CUE bottom | `managed-record` |
| Typed structured maps | [#7287](https://github.com/kubevela/kubevela/issues/7287) / [#7314](https://github.com/kubevela/kubevela/pull/7314) | `Map.OfObject` and `OfSchemaRef` describe structured values under dynamic keys without raw schema strings | `file-share` |
| Map-to-list transforms | [#7288](https://github.com/kubevela/kubevela/issues/7288) / [#7330](https://github.com/kubevela/kubevela/pull/7330) | Map entries can become list items while retaining both key and value | `file-share` |
| Actionable validation | [#7289](https://github.com/kubevela/kubevela/issues/7289) / [#7348](https://github.com/kubevela/kubevela/pull/7348) | Validator messages can interpolate the value being rejected | `zone-replica` |
| Multi-output health | [#7290](https://github.com/kubevela/kubevela/issues/7290) / [#7332](https://github.com/kubevela/kubevela/pull/7332) | Health expressions can target named outputs and aggregate output groups | `composite-store` |
| Fluent conditions | [#7353](https://github.com/kubevela/kubevela/issues/7353) / [#7352](https://github.com/kubevela/kubevela/pull/7352) | `NotMatches` emits CUE's native negative-regex operator | `tenant-space` |

Seven small components cover the eight merged capabilities. `file-share` combines
the two complementary collection APIs: a typed map-of-objects parameter and a
map-to-list workload transform.

The remaining four items are active roadmap work:

| Issue | Capability | Current state |
| --- | --- | --- |
| [#7281](https://github.com/kubevela/kubevela/issues/7281) | Structured values in `Lit` | Open; [PR #7313](https://github.com/kubevela/kubevela/pull/7313) is under review |
| [#7285](https://github.com/kubevela/kubevela/issues/7285) | Computed regex patterns | Open |
| [#7286](https://github.com/kubevela/kubevela/issues/7286) | Generic CUE standard-library calls | Open |
| [#7291](https://github.com/kubevela/kubevela/issues/7291) | Generated-CUE validation and evaluation | Open; [PR #7321](https://github.com/kubevela/kubevela/pull/7321) and [PR #7379](https://github.com/kubevela/kubevela/pull/7379) are under review |

## Presentation browser flow

The community-call script does not require slides. It starts in the browser with:

1. the official [DefKit overview](https://kubevela.io/docs/platform-engineers/defkit/overview/);
2. the [filtered twelve-issue list](https://github.com/kubevela/kubevela/issues?q=is%3Aissue%20author%3Akrishnankm%20Defkit%20in%3Atitle%20sort%3Acreated-asc);
3. this README's complete merged and active status tables;
4. [issue #7287](https://github.com/kubevela/kubevela/issues/7287) and merged [PR #7314](https://github.com/kubevela/kubevela/pull/7314); and
5. [issue #7283](https://github.com/kubevela/kubevela/issues/7283) and merged [PR #7339](https://github.com/kubevela/kubevela/pull/7339).

The presenter then switches once to the terminal demo. See
[`PRESENTATION.md`](PRESENTATION.md) for the word-for-word script and
[`PRESENTER-CHEAT-SHEET.md`](PRESENTER-CHEAT-SHEET.md) for the condensed cues.

## Quick start

Requirements:

- Go 1.23.8 or later;
- a `vela` CLI with `def vet` and offline `dry-run` support; and
- network access during the first `go mod download`.

The Go module pins the KubeVela merge revision that contains all eight merged
capabilities. No Kubernetes cluster is needed.

```bash
go mod download
./scripts/verify-demo.sh
```

The verification script:

1. generates all seven ComponentDefinitions;
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
  -f examples/route-catalog.yaml

vela dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml

vela dry-run --offline \
  -d generated/component \
  -f examples/composite-store.yaml
```

The intentionally invalid replica command is expected to fail. For that part of
the demo, focus on this output:

```text
parameter.replicas.1._validateReplicaZone."zone 'us-east-1' must not match the primary zone": conflicting values true and false
```

Use [`PRESENTER-CHEAT-SHEET.md`](PRESENTER-CHEAT-SHEET.md) for the recommended
flow. [`PRESENTATION.md`](PRESENTATION.md) contains the fuller talk track, while
[`DEMO-RUNBOOK.md`](DEMO-RUNBOOK.md) has the exact commands, handoff cues, and
recovery steps.

## Repository layout

```text
components/            Go-authored DefKit ComponentDefinitions
examples/              Applications used by offline dry-run
cmd/generate/           writes registered definitions as CUE
cmd/register/           emits the DefKit registry as JSON
scripts/verify-demo.sh  repeatable pre-call verification
PRESENTER-CHEAT-SHEET.md condensed run of show and speaker cues
PRESENTATION.md         timed narrative for the community call
DEMO-RUNBOOK.md         live commands and expected signals
```

## Scope of the demo

The examples emit resources under `example.com/v1alpha1`. They intentionally do
not require CRDs or controllers because the demo is about the definition-authoring
path: Go source, generated CUE, rendered resources, validation feedback, and
health-policy structure. A live cluster would add setup and failure modes without
making those DefKit behaviours clearer.
