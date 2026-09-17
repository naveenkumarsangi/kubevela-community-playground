# DefKit enhancements

This playground demonstrates eight merged DefKit improvements. Each section
explains when component authors might need the improvement, which API to use,
and how to try it in an example.

## Authoring structured resources

### Nested fields inside generated items

**Use it when:** a list comprehension must write nested fields, including
conditional siblings under the same parent.

```go
item.Set("spec.queueName", item.Var().Field("name"))
item.IfSet("maxReadOps", func() {
    item.Set("spec.throughput.maxReadOps", item.Var().Field("maxReadOps"))
})
```

Try it:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/queue-set.yaml
```

Look for nested `throughput` and `retention` fields under each queue.

- Source: [`components/queue_set.go`](../components/queue_set.go)
- Issue: [kubevela/kubevela#7282](https://github.com/kubevela/kubevela/issues/7282)
- Merged PR: [kubevela/kubevela#7302](https://github.com/kubevela/kubevela/pull/7302)

### Structured values in a map comprehension

**Use it when:** each source map entry must become an object under the original
key.

```go
body := defkit.NewResource("", "").
    Set("service", defkit.Reference("backendValue")).
    Set("metadata.sourceRoute", defkit.Reference("routeKey")).
    Ops()

catalog := defkit.ForEachMap().
    Over("parameter.routes").
    WithVars("routeKey", "backendValue").
    WithBody(body...)
```

Try it:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/route-catalog.yaml
```

A scalar input such as `checkout: checkout-api` becomes a structured
`routes.checkout` object.

- Source: [`components/route_catalog.go`](../components/route_catalog.go)
- Issue: [kubevela/kubevela#7283](https://github.com/kubevela/kubevela/issues/7283)
- Merged PR: [kubevela/kubevela#7339](https://github.com/kubevela/kubevela/pull/7339)

### Typed objects under dynamic map keys

**Use it when:** a map has arbitrary keys and every value follows the same object
schema.

```go
mountPoints := defkit.Map("mountPoints").
    Optional().
    OfObject(
        defkit.String("path").Required(),
        defkit.String("permissions").Optional(),
    )
```

With `OfObject`, each map value remains part of the typed builder model,
including its field descriptions, defaults, required markers, validators, and
import discovery. Use `OfSchemaRef` when the values share a reusable schema
definition.

- Source: [`components/file_share.go`](../components/file_share.go)
- Issue: [kubevela/kubevela#7287](https://github.com/kubevela/kubevela/issues/7287)
- Merged PR: [kubevela/kubevela#7314](https://github.com/kubevela/kubevela/pull/7314)

### Map-to-list transformation with key and value

**Use it when:** configuration is easiest to maintain as a keyed map, while the
workload API expects a list.

```go
mounts := defkit.NewArray().ForEachMapWithGuarded(
    mountPoints.IsSet(),
    mountPoints,
    func(entry *defkit.MapEntryBuilder) {
        entry.Set("name", entry.Key())
        entry.Set("path", entry.Value().Field("path"))
    },
)
```

Try it:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

The input keys become list-item names. Because map iteration order is not
guaranteed, do not rely on the order of the generated items.

- Source: [`components/file_share.go`](../components/file_share.go)
- Issue: [kubevela/kubevela#7288](https://github.com/kubevela/kubevela/issues/7288)
- Merged PR: [kubevela/kubevela#7330](https://github.com/kubevela/kubevela/pull/7330)

### Direct negative regex conditions

**Use it when:** a validator or conditional needs to match values that do not
satisfy a regex rule.

```go
defkit.LocalField("tenantName").NotMatches(`^[a-z0-9-]+$`)
tenantName.NotMatches(`^cust-`)
```

`NotMatches` emits CUE's native `!~` operator. It is available on both
`LocalFieldRef` and `StringParam`.

Try it:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/tenant-space.yaml
```

The example tenant does not start with `cust-`, so the rendered tier is
`internal`.

- Source: [`components/tenant_space.go`](../components/tenant_space.go)
- Issue: [kubevela/kubevela#7353](https://github.com/kubevela/kubevela/issues/7353)
- Merged PR: [kubevela/kubevela#7352](https://github.com/kubevela/kubevela/pull/7352)

## Validation and lifecycle behavior

### Health while status is absent

**Use it when:** a resource needs time to publish `status` or
`status.conditions` after creation.

```go
h := defkit.Health()

HealthPolicyExpr(h.And(
    h.Exists("status"),
    h.Exists("status.conditions"),
    h.AllTrue("Ready", "Synced"),
))
```

When `status` is missing, the health check evaluates the resource as not ready
instead of producing CUE bottom, an error value. Condition entries without
`type` are safely ignored.

- Source: [`components/managed_record.go`](../components/managed_record.go)
- Issue: [kubevela/kubevela#7284](https://github.com/kubevela/kubevela/issues/7284)
- Merged PR: [kubevela/kubevela#7299](https://github.com/kubevela/kubevela/pull/7299)

### Messages that identify the rejected value

**Use it when:** a validator runs over repeated objects and the user needs to know
which value violated the rule.

```go
defkit.ValidateValue(defkit.Interpolation(
    defkit.Lit("zone '"),
    defkit.LocalField("zone"),
    defkit.Lit("' must not match the primary zone"),
)).FailWhen(
    defkit.Eq(
        defkit.LocalField("zone"),
        defkit.Reference("parameter.primaryZone"),
    ),
)
```

Try the intentionally invalid input:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml
```

The command exits non-zero and includes:

```text
zone 'us-east-1' must not match the primary zone
```

- Source: [`components/zone_replica.go`](../components/zone_replica.go)
- Issue: [kubevela/kubevela#7289](https://github.com/kubevela/kubevela/issues/7289)
- Merged PR: [kubevela/kubevela#7348](https://github.com/kubevela/kubevela/pull/7348)

### Health across every resource in a component

**Use it when:** a component emits a primary resource, named auxiliary outputs,
and a repeated output group.

```go
h := defkit.Health()
vela := defkit.VelaCtx()

HealthPolicyExpr(h.And(
    h.At(vela.Output()).Condition("Ready").IsTrue(),
    h.At(vela.Outputs("accessPolicy")).Condition("Ready").IsTrue(),
    h.Every(
        defkit.OutputsWithPrefix("accessPoint"),
        func(item *defkit.HealthScope) defkit.HealthExpression {
            return item.Condition("Ready").IsTrue()
        },
    ),
))
```

`At` evaluates a health expression against a different resource. `Every`
evaluates an expression against each output selected by a prefix. Add
`AllowEmpty` when an empty output group should count as healthy.

Try it:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/composite-store.yaml
```

The result contains one `Store`, one `AccessPolicy`, and two `AccessPoint`
resources.

- Source: [`components/composite_store.go`](../components/composite_store.go)
- Issue: [kubevela/kubevela#7290](https://github.com/kubevela/kubevela/issues/7290)
- Merged PR: [kubevela/kubevela#7332](https://github.com/kubevela/kubevela/pull/7332)

## Related work in progress

Four related areas are still open or under review:

| Capability | Issue | Related PRs |
| --- | --- | --- |
| Structured values passed to `Lit` | [#7281](https://github.com/kubevela/kubevela/issues/7281) | [#7313](https://github.com/kubevela/kubevela/pull/7313) |
| Computed regex patterns | [#7285](https://github.com/kubevela/kubevela/issues/7285) | — |
| Generic CUE standard-library calls | [#7286](https://github.com/kubevela/kubevela/issues/7286) | — |
| Generated-CUE validation and evaluation helpers | [#7291](https://github.com/kubevela/kubevela/issues/7291) | [#7321](https://github.com/kubevela/kubevela/pull/7321), [#7379](https://github.com/kubevela/kubevela/pull/7379) |
