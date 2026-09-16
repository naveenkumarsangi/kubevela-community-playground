# Defkit feature examples

Each definition in `components/` demonstrates one defkit feature from a closed
kubevela issue, and each file here is an Application that exercises it. The file
name carries the issue number.

## Issue map

| Issue | Feature | Definition | Example |
| --- | --- | --- | --- |
| [#7282](https://github.com/kubevela/kubevela/issues/7282) | Nested field paths in `ItemBuilder` | `components/issue7282_queue_set.go` | `queue-set.yaml` |
| [#7284](https://github.com/kubevela/kubevela/issues/7284) | Health conditions that survive an absent status | `components/issue7284_managed_record.go` | `managed-record.yaml` |
| [#7288](https://github.com/kubevela/kubevela/issues/7288) | Guarded map-to-list comprehensions | `components/issue7288_file_share.go` | `file-share.yaml` |
| [#7289](https://github.com/kubevela/kubevela/issues/7289) | Expression-based validator messages | `components/issue7289_zone_replica.go` | `zone-replica.yaml` |
| [#7290](https://github.com/kubevela/kubevela/issues/7290) | Health expressions scoped to auxiliary outputs | `components/issue7290_composite_store.go` | `composite-store.yaml` |
| [#7353](https://github.com/kubevela/kubevela/issues/7353) | Negative regex conditions (`NotMatches`) | `components/issue7353_tenant_space.go` | `tenant-space.yaml` |

## What each example shows

### #7282 -- nested paths in `ItemBuilder` (`queue-set`)

`Resource.Set` took a nested path; `ItemBuilder.Set` treated the whole string as
one label, so `item.Set("throughput.maxReadOps", ...)` generated
`throughput.maxReadOps: ...` and `ToCue()` returned it without complaint. You
found out later, from `cue vet`, with `missing ',' in struct literal`.

`ItemBuilder.Set` now expands the path. The example covers a plain nested Set,
nested Sets inside `IfSet`/`IfNotSet`, and sibling writes under one parent
unifying into a single struct:

```cue
if v.maxReadOps != _|_ {
    spec: throughput: maxReadOps: v.maxReadOps
}
if v.maxWriteOps != _|_ {
    spec: throughput: maxWriteOps: v.maxWriteOps
}
```

### #7284 -- null-safe health conditions (`managed-record`)

A freshly applied resource has no status for the first few reconciles. The
condition comprehensions were emitted at the top level and dereferenced
`context.output.status.conditions` unconditionally, so CUE propagated bottom
through the `&&` chain instead of short-circuiting: the policy errored rather
than reporting "not healthy yet", and every reconcile logged a health-check
failure for the whole creation window.

The generated comprehensions now default a missing source and skip entries with
no `type`:

```cue
_readyCond: [ for c in *context.output.status.conditions | [] if c.type != _|_ if c.type == "Ready" { c } ]
```

Evaluated against the five fixtures from the issue -- `{}`, `{status: {}}`,
both-true, one-false, and a typeless entry alongside healthy ones -- `isHealth`
is `false, false, true, false, true`. No fixture produces bottom.

### #7288 -- map-to-list comprehensions (`file-share`)

Array-to-array worked and map-to-map worked; map-to-list with the key still in
scope did not. `ForEachWithVar` binds one iteration variable, so it can reach the
value but never the key, and `ForEachMap` has both variables but always emits a
struct.

`ForEachMapWithGuarded` gives both variables, a list result, and a guard for an
unset source. `ForEachMap` is unchanged.

```cue
mountPoints: [
    if parameter["mountPoints"] != _|_ for k, v in parameter.mountPoints {
        name: k
        path: v.path
        ...
    },
]
```

### #7289 -- expression-based validator messages (`zone-replica`)

`Validate` takes a Go string and the generator writes that same quoted string as
the field label in both branches, so a message can describe the rule but not the
value that broke it.

`ValidateValue` takes a `Value`. The generator binds it to a `let _message` and
uses `(_message)` as the computed key in both branches -- going through the
binding is what guarantees the two labels unify into one field rather than
becoming two:

```cue
_validateReplicaZone: {
    let _message = "zone '\(zone)' must not match the primary zone"
    (_message): true
    if zone == parameter.primaryZone {
        (_message): false
    }
}
```

Fixed messages stay the default; the example keeps one of each for contrast.

Still open, and split out of the issue deliberately: `ArrayParam.Validators`
emits `[...{...}]`, so no iterator is in scope and the message cannot name the
array index. It can name the field value.

### #7290 -- scoped and aggregated health (`composite-store`)

Health expressions were rooted at `context.output` with no way to point the same
builders at `context.outputs.<name>`, which left a component that creates several
resources together with raw CUE as its only complete option.

`At(ref)` roots an expression at any context reference, with the primary output
still the default. `Every(OutputsWithPrefix(p), fn)` requires every output whose
name starts with `p` to satisfy `fn`, and an empty match set reads unhealthy
rather than vacuously true -- `AllowEmpty()` opts into the other behaviour for a
genuinely optional group. That covers the three cases the issue named: a missing
target, absent `status.conditions`, and an empty collection all read unhealthy.

```cue
_accessPointItems: [ for k, v in (*context.outputs | {}) if k =~ "^accessPoint" { v } ]
```

### #7353 -- negative regex conditions (`tenant-space`)

Defkit had typed helpers for positive regex matching but nothing symmetric for
the negative case, even though CUE has a native `!~`. Callers wrapped a positive
match in `Not(...)`, which is valid but less discoverable and emits
`!(x =~ "p")`.

`NotMatches` now exists on `StringParam` (runtime conditions) and
`LocalFieldRef` (validators), and emits `!~` directly:

```cue
if parameter.tenantName !~ "^cust-" {
    tier: "internal"
}
```

## Seeing the generated CUE

This needs no cluster:

```bash
go run ./cmd/generate generated                     # CUE for every definition
vela def vet generated/component/queue-set.cue      # validates offline
```

All six definitions pass `vela def vet`. With no cluster reachable it logs an
error about loading external cuex packages first and then reports
`Validation ... succeed.` -- the validation itself does not need the cluster.

