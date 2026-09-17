# DefKit community-call cheat sheet

Target: 8–9 minutes. No slides. Share the browser first, then switch once to the
Demo terminal.

Read the blockquoted paragraphs nearly verbatim. Headings, commands, and plain
instructions are stage directions and should not be spoken.

## Browser tabs, left to right

1. https://kubevela.io/docs/platform-engineers/defkit/overview/
2. https://github.com/kubevela/kubevela/issues?q=is%3Aissue%20author%3Akrishnankm%20Defkit%20in%3Atitle%20sort%3Acreated-asc
3. https://github.com/naveenkumarsangi/kubevela-community-playground/tree/demo-2026-09-17
4. https://github.com/kubevela/kubevela/issues/7287
5. https://github.com/kubevela/kubevela/pull/7314
6. https://github.com/kubevela/kubevela/issues/7283
7. https://github.com/kubevela/kubevela/pull/7339

## 0:00–0:55 — Official DefKit docs

Point to the overview and Benefits section.

> I'll start with the tool behind these examples. DefKit is KubeVela's Go SDK
> for authoring X-Definitions. With it, we build the schema, template, validation,
> and health policy in Go, and DefKit generates the CUE that KubeVela already
> understands.
>
> The benefits that mattered for us were normal Go tooling, IDE support, testing,
> package distribution, and compatibility with the existing CUE runtime.

Do not read every benefit.

## 0:55–1:30 — The twelve issues

Show the filtered issue list.

> We used DefKit in two ways: we created new OAM components, and we moved selected
> definitions from hand-written CUE to Go. As we tried it across more component
> shapes, a few patterns kept coming up that were useful beyond any one component.
>
> We captured those cases as focused upstream issues. There are twelve in total:
> eight are merged, and four remain open or under review. They cover typed
> authoring patterns, validation and health behaviour, and multi-resource
> lifecycle support.

Do not read all twelve titles.

## 1:30–2:00 — The complete merged set

Show the **Merged outcomes demonstrated here** table in the playground README.

> Before I open individual examples, this table shows the complete status. On
> the authoring side, the merged work covers nested paths, structured map schemas
> and bodies, map-to-list transforms, and negative-regex conditions. On the
> feedback and lifecycle side, it adds health before status exists,
> value-specific validation, and multi-output health.
>
> We are not going to open eight pull requests in a short demo. I'll focus on two
> recent merges, and the terminal flow will exercise five representative
> capabilities. The repository has runnable examples for all eight, and the four
> active items are listed below the table.

Point to the eight merged rows and four active items. Do not read every row.

## 2:00–2:35 — #7287 and PR #7314

Issue:

> One pattern we kept running into was a dynamic-key map where each value is an
> object. A structured value schema still required a raw CUE schema string, so
> issue #7287 captured that reusable shape.

Merged PR:

> PR #7314 adds typed `OfObject` and `OfSchemaRef` builders. A community
> contributor implemented the change, and we verified it against the component
> shape that originally motivated the issue.

Point to **Merged** and the API example. Do not open the diff.

## 2:35–3:05 — #7283 and PR #7339

Issue:

> `ForEachMap.WithBody` described a structured object for each map entry, but
> those body operations were not appearing in the generated output.

Merged PR:

> PR #7339 makes the generator render the body under every original key and
> collect its required imports. We contributed the fix and verified it with both
> compiled and evaluated CUE.

Point to **Merged** and the before/after shape.

## 3:05–3:15 — Terminal handoff

> Those are two recent examples from eight merged improvements. I'll switch to
> the terminal now and show how they fit into a normal DefKit authoring flow.

Switch once to the Demo terminal. Do not return to the browser during the demo.

## 3:15–3:40 — Generate

```bash
rm -rf generated
go run ./cmd/generate generated
```

> Seven Go-authored definitions cover the eight merged capabilities. The
> generator produces standard CUE definitions for the existing KubeVela runtime.

## 3:40–4:55 — Typed collections

```bash
bat --paging=never --line-range 37:75 \
  components/issues7287_7288_file_share.go
```

> Looking at the source, `OfObject` defines the object under each dynamic map key.
> The workload expects a list, so `ForEachMapWithGuarded` carries both key and
> value forward while building each item.

```bash
vela dry-run --offline -d generated/component \
  -f examples/file-share.yaml
```

Point to `cache`, `reports`, and the default `0755` permission.

```bash
vela dry-run --offline -d generated/component \
  -f examples/route-catalog.yaml
```

> Here the result stays map-shaped, but every backend becomes a structured object
> under its original route key. That's the merged `WithBody` behaviour.

## 4:55–5:35 — Validation feedback

```bash
vela dry-run --offline -d generated/component \
  -f examples/zone-replica-invalid.yaml
```

> We expect this rejection because a replica is using the primary zone. The useful
> part is that the message identifies `us-east-1`, the exact value to correct.

Do not explain the `conflicting values` suffix.

## 5:35–6:55 — Multi-resource health

```bash
vela dry-run --offline -d generated/component \
  -f examples/composite-store.yaml
```

Point to `Store/orders`, `AccessPolicy/orders-policy`, and both access points.

```bash
bat --paging=never --line-range 51:68 \
  components/issue7290_composite_store.go
```

> For a multi-resource component, generation is only half the job. This typed
> health policy checks the primary output, the named policy output, and every
> output with the `accessPoint` prefix.

## 6:55–7:50 — Summary

> In the demo, we used five of the eight merged capabilities. The repository has
> runnable examples for nested item paths, health before status exists, and
> negative-regex conditions as well.
>
> Four areas remain active: structured literals, computed regex patterns, generic
> standard-library calls, and generated-CUE validation helpers.
>
> We used DefKit for both new components and migrations from CUE, and that work
> produced reusable improvements for the wider KubeVela community.

## 7:50–8:10 — Close

> Thank you to everyone who reviewed, implemented, and verified these changes.
> The playground repository has all of the runnable examples and the complete
> issue status.

Stop. Leave time for questions.

## Recovery

```bash
bat --paging=never /tmp/defkit-demo-backup/file-share.txt
bat --paging=never /tmp/defkit-demo-backup/route-catalog.txt
bat --paging=never /tmp/defkit-demo-backup/validation.txt
bat --paging=never /tmp/defkit-demo-backup/composite-store.txt
```
