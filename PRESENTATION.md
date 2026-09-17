# Community-call script: DefKit lessons from real component work

Target: 8–9 minutes. No slides. The presentation moves from the official DefKit
documentation, to the upstream issues and merged pull requests, and then to the
terminal demo.

## Pages to open before the call

Open these browser tabs in this exact order:

1. [DefKit overview](https://kubevela.io/docs/platform-engineers/defkit/overview/)
2. [DefKit issues authored by `krishnankm`](https://github.com/kubevela/kubevela/issues?q=is%3Aissue%20author%3Akrishnankm%20Defkit%20in%3Atitle%20sort%3Acreated-asc)
3. [Playground status table](https://github.com/naveenkumarsangi/kubevela-community-playground/tree/demo-2026-09-17)
4. [Issue #7287: structured map values](https://github.com/kubevela/kubevela/issues/7287)
5. [PR #7314: `OfObject` and `OfSchemaRef`](https://github.com/kubevela/kubevela/pull/7314)
6. [Issue #7283: `ForEachMap.WithBody`](https://github.com/kubevela/kubevela/issues/7283)
7. [PR #7339: render `ForEachMap` body operations](https://github.com/kubevela/kubevela/pull/7339)

Keep the Demo terminal open separately. Share the browser first, then switch once
to the terminal. Do not move repeatedly between browser and terminal.

## 0:00–0:55 — Explain DefKit from the official documentation

**Screen:** DefKit overview. Keep the page near the title, then scroll to
**Benefits**.

**Say:**

> I'll start with the tool behind these examples. DefKit is KubeVela's Go SDK
> for authoring X-Definitions. With it, we build the component schema, resource
> template, validation, and health policy in Go, and DefKit generates the CUE
> that KubeVela already understands.
>
> The benefits that mattered for this work are the ones you can see here: normal
> Go tooling, IDE support, unit testing, package-based distribution, and
> compatibility with the existing CUE runtime. CUE is still the runtime contract;
> DefKit gives us a maintainable authoring layer in Go.

Do not read every benefit from the page. Point to the list and name only the
benefits relevant to the component work.

## 0:55–1:30 — Show the real contribution set

**Screen:** filtered GitHub issue list.

**Say:**

> We used DefKit in two ways: we created new OAM components, and we moved selected
> existing definitions from hand-written CUE to Go. As we tried it across more
> component shapes, a few patterns kept coming up that were useful beyond any one
> component.
>
> We captured those cases as focused upstream issues. You're looking at the full
> issue set here, including both open and closed items. There are twelve in total:
> eight are now merged, and four remain open or under review. They group into
> three broad areas—typed authoring patterns, validation and health behaviour, and
> multi-resource lifecycle support.

On the issue list, point out that both open and closed items are present. Do not
read twelve titles or open every issue.

## 1:30–2:00 — Acknowledge the complete merged set

**Screen:** playground README, at **Merged outcomes demonstrated here**.

**Say:**

> Before I open individual examples, this table shows the complete status. On
> the authoring side, the merged work covers nested paths, structured map schemas
> and bodies, map-to-list transforms, and negative-regex conditions. On the
> feedback and lifecycle side, it adds health before status exists,
> value-specific validation, and multi-output health.
>
> We are not going to open eight pull requests in a short demo. I'll focus on two
> of the most recent merges, and the terminal flow will exercise five
> representative capabilities. The repository has runnable examples for all
> eight, and the four active items are listed directly below this table.

Point to the eight-row merged table and the four-item active-work table. Do not
read each row aloud.

## 2:00–2:35 — Show newly merged structured-map support

**Screen:** issue #7287, followed by PR #7314.

**Say on the issue:**

> One pattern we kept running into was a dynamic-key map where each value is an
> object. We could already describe a simple map, but a structured value schema
> still required a raw CUE schema string. Issue #7287 captured that shape so it
> could be supported as a reusable pattern.

**Switch to PR #7314 and say:**

> Over in PR #7314, you can see that the change is now merged. It adds typed
> `OfObject` and `OfSchemaRef` builders, which means the complete map value
> schema—including required, optional, and defaulted fields—can stay in Go. A
> community contributor implemented the change, and we verified it against the
> component shape that originally motivated the issue.

Point to the green **Merged** state and the API example. Do not inspect the file
diff during the call.

## 2:35–3:05 — Show newly merged map-body generation

**Screen:** issue #7283, followed by PR #7339.

**Say on the issue:**

> The second map pattern used `ForEachMap.WithBody` to build a structured object
> for each input entry. The authoring API was already there, but those body
> operations were not appearing in the generated output. Issue #7283 isolated
> that specific behaviour.

**Switch to PR #7339 and say:**

> PR #7339 is merged as well. The generator now renders the body under each
> original map key and collects any imports required by the body values. We
> contributed this fix and verified it with both compiled and evaluated CUE.

Again, point to **Merged** and the before-and-after generated shape. Do not read
the full PR description.

## 3:05–3:15 — Move once to the terminal

**Say:**

> Those are two recent examples from the eight merged improvements. Rather than
> opening every pull request, I'll switch to the terminal now and show how these
> capabilities fit into a normal DefKit authoring and rendering flow.

Switch to the Demo terminal. Do not return to the browser unless answering a
question after the demo.

## 3:15–3:40 — Generate the definitions

**Run:**

```bash
rm -rf generated
go run ./cmd/generate generated
```

**Say:**

> In this repository, seven small Go-authored ComponentDefinitions cover the
> eight merged capabilities. When I run the generator, the output is a set of
> standard CUE definitions, so the KubeVela runtime does not need a separate
> DefKit execution path.

## 3:40–4:55 — Demonstrate typed collection authoring

**Show the Go source:**

```bash
bat --paging=never \
  --line-range 37:75 \
  components/issues7287_7288_file_share.go
```

**Say:**

> Looking at the source, `OfObject` defines the object that sits under each
> dynamic map key. The workload expects a list, so `ForEachMapWithGuarded` carries
> both the key and the value forward while it builds each list item. Together,
> those two APIs cover the full path from a typed map-of-objects input to a
> list-shaped workload field.

**Render it:**

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

**Point to `spec.mountPoints` and say:**

> In the rendered `spec.mountPoints`, the original map keys have become the
> `cache` and `reports` item names. The first entry also picked up the default
> permission because the Application did not provide one.

**Render the `WithBody` example:**

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/route-catalog.yaml
```

**Say:**

> In this example, the result stays map-shaped, but each scalar backend becomes a
> structured object under its original route key. That's the merged `WithBody`
> behaviour we just looked at in the pull request.

## 4:55–5:35 — Demonstrate value-specific validation

**Run:**

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml
```

**Say:**

> We expect this one to be rejected because a replica is using the primary zone.
> The useful part is that the generated message includes `us-east-1`, the exact
> value that needs to be corrected. The user can go straight to that entry
> instead of searching the list for the one that broke the rule.

Point only to `zone 'us-east-1' must not match the primary zone`. Do not explain
CUE's `conflicting values` suffix.

## 5:35–6:55 — Demonstrate multi-resource health

**Run:**

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/composite-store.yaml
```

**Point out:**

- primary `Store/orders`;
- named `AccessPolicy/orders-policy`;
- repeated `AccessPoint/orders-ap-1`; and
- repeated `AccessPoint/orders-ap-2`.

**Show the readable Go health policy:**

```bash
bat --paging=never \
  --line-range 51:68 \
  components/issue7290_composite_store.go
```

**Say:**

> For a multi-resource component, generating the resources is only half the job;
> the health policy has to cover all of them as well. In this typed policy, we
> check the primary output, the separately named access policy, and every output
> whose name has the `accessPoint` prefix.
>
> So DefKit generates all four resources, and the same Go authoring model also
> describes the lifecycle of the complete component.

## 6:55–7:50 — Summarize the result

**Keep the terminal visible. Say:**

> In that live flow, we used five of the eight merged capabilities: typed
> structured maps, map-to-list transformation, structured `ForEachMap` bodies,
> value-specific validation, and multi-output health.
>
> The repository has runnable examples for the other areas too, including nested
> item paths, health evaluation before status exists, and native negative-regex
> conditions. Of the twelve original items, four remain active: structured
> literals, computed regex patterns, generic standard-library calls, and
> generated-CUE validation helpers.
>
> We were able to use DefKit for both new OAM components and migrations from
> hand-written CUE. Along the way, that component work also produced reusable
> improvements for the wider KubeVela community.

## 7:50–8:10 — Close

**Say:**

> Thank you to everyone who reviewed, implemented, and verified these changes.
> The playground repository has all of the examples and the complete issue
> status, and we're happy to dig into any of the individual designs.

Stop. Leave the remaining time for questions.

## Likely questions

**Why not demonstrate every merged issue?**
Walking through eight separate examples would make the overall workflow harder
to see. The live demo focuses on five representative capabilities, and the
repository has runnable examples for all eight.

**Does DefKit replace CUE?**
No. CUE remains the runtime contract. DefKit is the Go authoring layer, and it
generates CUE for the existing KubeVela runtime.

**Why is the demo offline?**
The contribution surface here is definition authoring, CUE generation,
validation, and rendering. An offline dry-run exercises all of those boundaries
without bringing unrelated CRD or controller setup into the demo.

**What remains open?**
Four areas remain open or under review: structured literal values, computed
regex patterns, generic standard-library calls, and generated-CUE
validation/evaluation helpers.
