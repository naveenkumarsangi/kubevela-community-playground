# Community-call talk track: DefKit lessons from real component work

Target length: 9–10 minutes, including the demo.

Roles below use **Presenter** for the narrative and **Demo partner** for the
terminal. If one person presents, keep the same order and ignore the handoff.

## The message to leave with the audience

We adopted DefKit to make OAM component authoring easier to maintain in Go. It
worked for both new components and migrations from hand-written CUE. Real
definitions highlighted opportunities to expand the typed builders; we turned
those cases into reproducible upstream proposals, contributed several changes,
and collaborated with other contributors on the rest. Six improvements are now
merged.

## 0:00–0:45 — Start with the goal, not the issue list

**Presenter**

> Our work started with a practical goal: build new OAM components and move a few
> existing components from hand-written CUE to DefKit. We wanted the definition
> codebase to feel like the rest of our Go platform code—easy to navigate, review,
> test, and evolve.
>
> As we moved through that work, we found patterns worth making reusable in
> DefKit. We reduced each one to a focused example, discussed it with the
> maintainers, and contributed the resulting enhancements upstream.

## 0:45–1:50 — Briefly explain DefKit and why it fit

**Presenter**

> DefKit is KubeVela's Go SDK for authoring X-Definitions. We write the schema,
> resource template, validation, and health policy with fluent Go builders, and
> DefKit generates the CUE that KubeVela already evaluates.
>
> In practice, that gave us familiar IDE support, reusable builders, normal Go
> tests and module versioning, and less context switching for engineers already
> working in Go. It also let us introduce new OAM components and migrate existing
> ones incrementally.
>
> CUE remains the generated contract and KubeVela's runtime engine. DefKit gives
> us a typed Go authoring layer while preserving raw CUE as an option for
> specialized cases.

## 1:50–3:00 — Explain how component work shaped the enhancements

**Presenter**

> As we applied DefKit to more demanding components, several reusable extension
> opportunities became clear.
>
> One was broader authoring expressiveness: nested writes inside generated list
> items, turning a keyed map into a list without losing the key, and saying that
> a value must not match a pattern.
>
> Another was richer feedback and lifecycle behaviour: validation messages could
> identify the exact value being rejected, and a resource with no status yet
> could be treated as normally not ready while it starts reconciling.
>
> Multi-resource health was a third opportunity. A component can emit a primary
> resource, a named policy resource, and a collection of auxiliary resources.
> Extending the health builders to those outputs lets the complete policy stay in
> the same typed Go authoring model.
>
> We turned these opportunities into focused examples and initially opened eleven
> upstream issues. A later migration case added one more focused enhancement for
> negative regex matching.

## 3:00–3:40 — Summarize the result without reading issue numbers

**Presenter**

> Six improvements from that effort are now merged. They make nested and
> collection-heavy templates easier to express, make validation feedback more
> specific, and let health policies describe the whole component rather than only
> its primary resource.
>
> Some fixes came from us and some were implemented by other KubeVela
> contributors. We also verified the changes against the component shapes that
> exposed them. That's the part of the community process we wanted to share: real
> use, a small reproducer, design review, implementation, and verification.

**Handoff**

> Instead of walking through six pull requests, we'll show one authoring flow and
> three outcomes that users can see directly.

## 3:40–7:40 — Demo

Follow [`DEMO-RUNBOOK.md`](DEMO-RUNBOOK.md).

### Demo beat 1: Go definitions become normal KubeVela definitions

**Demo partner**

Generate the definitions.

> These six definitions are ordinary Go packages using DefKit. Generating them
> produces standard CUE ComponentDefinitions. There is no special runtime and no
> cluster dependency in this demo.

Show `components/issue7288_file_share.go`, then dry-run `file-share.yaml`.

> This component accepts mount points as a map because names make configuration
> easier to manage. The target API wants a list. The merged map-to-list builder
> keeps both the map key and value, so the rendered resource contains named list
> entries and applies a default permission where the input omitted one. Before
> this was added, this shape needed raw CUE or awkward restructuring.

### Demo beat 2: validation identifies the actual bad value

Run the intentionally invalid zone-replica Application.

> The second replica uses the primary zone. The command is expected to fail, but
> notice the message: it includes `us-east-1`, the value that violated the rule.
> With only a fixed string, users knew the rule, but they still had to search the
> input for the offending item. An expression-based message points to the bad
> value right in the validation error.

### Demo beat 3: one component can own several resources

Dry-run `composite-store.yaml`, then show the health-policy excerpt in the
generated definition.

> The component renders a primary store, a named access policy, and two access
> points. Its health policy uses the same typed builder for the primary output, a
> named output, and every output with the `accessPoint` prefix.
>
> This mattered for migrations. Multi-resource definitions were one of the places
> where we could generate the resources in Go but still had to drop into raw CUE
> to describe health. The merged scoping and aggregation APIs close that gap.

Optional, if there is time:

> The repository also contains the other merged cases: nested item paths,
> null-safe health before status exists, and a fluent `NotMatches` condition.
> They all run through the same generate, vet, and dry-run workflow.

## 7:40–8:45 — Connect the demo back to the effort

**Presenter**

> The examples are intentionally small, but they came from larger component
> work. We used DefKit to build new OAM components and migrate selected CUE
> components to the Go authoring model. Each case we moved from raw CUE into typed
> builders made the definition easier for a Go-focused team to review and change.
>
> We kept each enhancement focused. We reproduced the concrete component shape,
> discussed the smallest reusable API, and tested the generated CUE—not just the
> Go string returned by the builder. That last point mattered because several
> behaviours only became visible when CUE was parsed or evaluated.

## 8:45–9:35 — Status and next steps

**Presenter**

> Six improvements are merged and demonstrated here. Other items from the
> original set remain open or under review, including broader structured-map
> support, generated-CUE evaluation helpers, and generic standard-library calls.
> We'll keep prioritizing them based on which ones remove real escape hatches
> from component work, rather than by issue number.
>
> The playground pins a post-merge KubeVela revision, and the whole demo runs
> offline after dependencies are downloaded. That makes the examples useful both
> for this call and as regression-sized references for future DefKit users.

## 9:35–10:00 — Close

**Presenter**

> The short version: we chose DefKit to make OAM definitions easier for a
> Go-based platform team to maintain. Building and migrating real components
> showed us several places where the typed API could support more component
> patterns. We turned those opportunities into focused contributions. The merged
> changes let DefKit express more, behave more predictably, and better support
> multi-resource components.
>
> Thanks to everyone who reviewed, implemented, and verified these changes.

## Questions we are likely to get

**Does DefKit replace CUE?**  
No. DefKit is an authoring SDK. It generates CUE, and KubeVela continues to use
CUE at runtime.

**When do you still use raw CUE?**  
Raw CUE remains a useful escape hatch for genuinely specialized cases. The
contributed APIs bring recurring component patterns into the typed authoring
model while preserving that flexibility.

**Were all fixes written by the same team?**  
No. We raised and verified the problem set, contributed several fixes, and other
community members implemented others. This was a community effort.

**Why is the demo offline instead of running on Kubernetes?**  
These changes affect definition authoring, generated CUE, rendering, validation,
and health-policy structure. Offline generation and dry-run exercise those
boundaries directly. Installing fake CRDs and controllers would add setup but no
useful evidence.

**Are all twelve issues complete?**  
No. Six demonstrated improvements are merged. The remaining work is still
prioritized by concrete component needs and maintainer feedback.
