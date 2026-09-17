# DefKit community-call demo runbook

The full presentation uses a three-minute browser walkthrough followed by about
five minutes in the terminal. The terminal section uses offline generation and
dry-run, so there is no cluster to prepare or debug during the call.

## Before the call

From the repository root:

```bash
go mod download
./scripts/verify-demo.sh
```

A successful preflight ends with:

```text
rejected with: zone 'us-east-1' must not match the primary zone
Demo verification passed.
```

Run it once the day before and once shortly before the call. Download the
dependencies in advance and leave the generated directory in place. Increase
the terminal font size and open this file in a second window.

If the `vela` binary is not on `PATH`, point the scripts and commands at it:

```bash
export VELA_BIN=/path/to/vela
```

## Browser lead-in

Follow the seven tabs in `PRESENTATION.md`: official DefKit overview, the
filtered twelve-issue list, the playground README status tables, issue #7287
with merged PR #7314, and issue #7283 with merged PR #7339. Spend about three
minutes in the browser.

After PR #7339, say:

> Those are two recent examples from eight merged improvements. Rather than
> opening every pull request, I'll switch to the terminal now and show how these
> capabilities fit into a normal DefKit authoring and rendering flow.

Switch once to the Demo terminal. Do not return to the browser during the demo.

## Terminal flow

### 1. Establish the authoring pipeline — 30 seconds

```bash
rm -rf generated
go run ./cmd/generate generated
```

Expected signal: seven `.cue` files are written under `generated/component/`.

Say:

> These seven examples cover eight merged capabilities. They are authored as
> normal Go packages, and DefKit turns them into standard CUE
> ComponentDefinitions without adding anything to the KubeVela runtime.

Do not open every file. The audience only needs to see that generation is one
repeatable step.

### 2. Show typed collection authoring — 110 seconds

Show the structured-map parameter and map-to-list builder:

```bash
bat --line-range 27:78 components/issues7287_7288_file_share.go
```

Then render the Application:

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

Point to the typed input in Go and the rendered `spec.mountPoints` list.

Say:

> `OfObject` describes the value under every dynamic map key with typed DefKit
> fields, so the parameter no longer needs a raw schema string. The workload
> expects a list, and the map-to-list builder keeps both the key and value while
> applying a default permission. Together these APIs cover the full path from a
> typed map-of-objects input to a list-shaped workload field.

The order of map entries is not part of the contract. Do not describe `cache`
as always appearing first.

Then show the second newly merged map transform:

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/route-catalog.yaml
```

Point to `spec.routes.checkout` and `spec.routes.search`.

Say:

> This case remains map-shaped, but each scalar input becomes a structured value
> under its original key. `ForEachMap.WithBody` already represented that intent;
> the generator now renders the body operations for every entry. These two
> examples cover both common structured-map directions without raw CUE.

### 3. Show more useful validation feedback — 60 seconds

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml
```

This intentionally invalid input is rejected. Point to the value-specific message:

```text
parameter.replicas.1._validateReplicaZone."zone 'us-east-1' must not match the primary zone": conflicting values true and false
```

Say:

> The generated validator includes the actual zone that violated the rule. A
> fixed message described the rule; this message tells the user what to fix
> immediately.

If someone asks where that message comes from:

```bash
bat --line-range 45:52 generated/component/zone-replica.cue
```

### 4. Show the multi-resource shape — 75 seconds

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/composite-store.yaml
```

Point out the four rendered resources:

- primary `Store`;
- named `AccessPolicy`;
- `accessPoint1`; and
- `accessPoint2`.

Then show the readable Go health policy:

```bash
bat --paging=never \
  --line-range 51:68 \
  components/issue7290_composite_store.go
```

Say:

> The typed policy checks the primary store, the separately named access policy,
> and every output whose name starts with `accessPoint`. DefKit is generating all
> four resources and describing the lifecycle of the complete component in the
> same Go authoring model.

### 5. Hand back — 20 seconds

Say:

> This flow showed five of the eight merged capabilities. The repository also
> covers nested item paths, an unhealthy result instead of an error when status
> is absent, and a fluent negative-regex condition. Every example uses the same
> generate, vet, and dry-run workflow.

Hand back to the presenter for the contribution summary and close.

## Optional 30-second add-on

If the previous steps finished early, show native negative matching:

```bash
"${VELA_BIN:-vela}" dry-run --offline \
  -d generated/component \
  -f examples/tenant-space.yaml

bat --line-range 24:34 generated/component/tenant-space.cue
```

The input does not start with `cust-`, so `NotMatches` renders `tier: internal`.
Keep this optional; it should not push out the multi-resource example.

## Recovery

### Generation fails while downloading modules

The live demo should not be the first module download. Run this before the call:

```bash
go mod download
go run ./cmd/generate generated
```

### `vela` is missing

```bash
export VELA_BIN=/absolute/path/to/vela
./scripts/verify-demo.sh
```

### The invalid example returns a non-zero exit code

That is the intended result. Continue once the error contains:

```text
zone 'us-east-1' must not match the primary zone
```

### Generated files are stale

```bash
rm -rf generated
go run ./cmd/generate generated
```

### A live command still fails

Do not spend the remaining presentation debugging. Show these verified signals:

1. `file-share` uses a typed map schema and renders a list with names retained
   from map keys;
2. `route-catalog` renders a structured object under every original map key;
3. the invalid replica identifies `us-east-1` in its validation message; and
4. `composite-store` renders one primary and three auxiliary resources.

Then continue with the contribution summary. The focus is the authoring process
and the merged capabilities, not the terminal session itself.
