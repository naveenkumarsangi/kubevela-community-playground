# DefKit demo: presenter cheat sheet

## The order to use

Use this order. It gives the audience the reason for the work before showing the
implementation, and the demo moves from the easiest result to the most advanced.

1. The component goal
2. Why DefKit fit that goal
3. What the real component work helped improve
4. Three demo outcomes
5. Community result and close

Aim to finish in eight to nine minutes. Keep the last minute as buffer rather
than filling the full ten minutes with prepared speech.

## 0:00–0:35 — Open with the component goal

Show: title slide only.

Say:

> We used DefKit for two related goals: creating new OAM components and moving a
> few existing definitions from hand-written CUE to a Go authoring model. Today
> we want to show what that experience enabled and how it led to reusable DefKit
> improvements upstream.

Do not open with issue numbers or pull requests. The audience first needs to know
what you were trying to deliver.

## 0:35–1:20 — Explain DefKit in one picture

Show:

```text
Go definition using DefKit
            ↓
generated CUE ComponentDefinition
            ↓
      normal KubeVela runtime
```

Say:

> DefKit is KubeVela's Go SDK for authoring X-Definitions. We describe the schema,
> template, validation, and health policy with typed Go builders. DefKit generates
> the CUE that KubeVela already understands.
>
> This gave us familiar IDE support, reusable builders, normal Go tests, and an
> incremental path for migrating existing definitions. CUE remains the runtime
> contract; DefKit gives us a maintainable Go authoring layer.

Do not explain the complete DefKit API. The audience only needs this mental model.

## 1:20–2:10 — Describe the effort and the opportunities

Show: one slide with two inputs and three outcome areas.

```text
New components + CUE migrations
                ↓
  ┌─────────────┼─────────────┐
  │             │             │
Authoring    Validation    Multi-resource
patterns     and health    lifecycle
```

Say:

> As we applied DefKit to more component shapes, we found patterns worth making
> reusable: typed schemas for structured maps, richer map transformations, more
> specific validation feedback, predictable health while resources start
> reconciling, and health policies that cover every resource emitted by a
> component.
>
> We turned those component cases into focused upstream proposals. Eight of the
> twelve improvements are now merged, with changes contributed both by us and by
> other KubeVela community members.

Do not describe all twelve issues. If someone asks, the full status is in the
repository and can be covered during questions.

## 2:10–2:20 — Handoff

Presenter:

> Instead of walking through pull requests, we will show three themes from one
> normal DefKit authoring flow.

Demo partner:

> I will start with typed collection authoring, then show validation feedback,
> and finally a multi-resource health policy.

## 2:20–3:00 — Establish the generation path

Run:

```bash
rm -rf generated
go run ./cmd/generate generated
```

Say:

> These seven Go-authored ComponentDefinitions cover eight merged capabilities.
> Generation produces standard CUE; there is no extra runtime and this part needs
> no cluster.

Show one source-to-CUE pair:

```bash
bat --line-range 27:78 components/issues7287_7288_file_share.go
bat --line-range 20:45 generated/component/file-share.cue
```

Do not open every source file or generated definition.

## 3:00–4:30 — Demo 1: typed collection authoring

Run:

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

Point to the `OfObject` schema in Go and the rendered `spec.mountPoints` list.

Say:

> `OfObject` describes the structured value under every dynamic map key without
> a raw schema string. The workload expects a list, so the map-to-list builder
> then keeps both the key and value and applies a default permission. This covers
> the complete path from typed input schema to rendered workload output.

Then run:

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/route-catalog.yaml
```

Say:

> This transform stays map-shaped but builds a structured object under every
> original key. The merged `ForEachMap.WithBody` support now renders those body
> operations for each entry.

The order of map entries is not guaranteed. Do not make the explanation depend
on either map appearing in a particular order.

## 4:30–5:15 — Demo 2: specific validation feedback

Run:

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml
```

Point to:

```text
zone 'us-east-1' must not match the primary zone
```

Say:

> This input is intentionally rejected because a replica uses the primary zone.
> The important change is the message: it includes the value that needs to be
> corrected. The user no longer has to search the list to work out which entry
> broke the rule.

The non-zero exit code is expected. Continue immediately after showing the
message.

## 5:15–6:30 — Demo 3: multi-resource health

Run:

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/composite-store.yaml
```

Point out:

- the primary `Store`;
- the named `AccessPolicy`;
- `accessPoint1`; and
- `accessPoint2`.

Then show:

```bash
bat --line-range 13:21 generated/component/composite-store.cue
```

Say:

> This component owns four resources, so health has to describe the whole group.
> The typed policy checks the primary output, the named policy output, and every
> output with the `accessPoint` prefix. This is especially useful for migrations:
> resource generation and health can now stay in the same Go authoring model.

Stop here. Do not add the negative-regex example unless the earlier steps finished
well ahead of time.

## 6:30–7:15 — Explain the contribution result

Presenter:

> These examples came from real component shapes, but we kept each upstream
> proposal focused and reusable. We reproduced the case, discussed the API with
> maintainers, tested the generated CUE, and verified the result against the
> component pattern that motivated it.
>
> Some changes were implemented by us and others by community contributors. That
> collaboration gave us eight merged improvements across authoring, validation,
> and multi-resource lifecycle handling.

## 7:15–7:45 — Close

Say:

> DefKit helped us build new OAM components and move selected CUE definitions to
> a Go-based authoring workflow. Applying it to real component work also gave us
> clear opportunities to contribute reusable improvements. The result is a more
> expressive and predictable authoring experience for the next set of components
> as well.

Then stop. Leave time for questions.

## What to show

- One DefKit Go component and its matching generated CUE excerpt.
- Two successful structured-map transformations.
- One value-specific validation message.
- One multi-resource render and its health-policy excerpt.
- A final slide saying eight of twelve improvements are merged through team and
  community collaboration.

## What not to show

- All twelve issue descriptions.
- Seven source files one by one.
- Pull-request diffs or test files.
- The complete generated CUE documents.
- k3d, CRD installation, or fake controllers.
- The optional negative-regex example unless there is extra time.

## If time is cut to five minutes

1. Give the goal and DefKit explanation in one minute.
2. Run the `file-share` and multi-resource examples.
3. Close with eight of twelve improvements merged through community collaboration.
