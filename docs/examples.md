# Example catalog

Before running any example, generate the definitions:

```bash
go run ./cmd/generate generated
```

Every command below runs offline against the generated definition directory.
The rendered resources use illustrative `example.com/v1alpha1` kinds and are not
intended to be applied to a cluster.

## Queue set

**Explore:** nested `ItemBuilder` paths and conditional sibling fields.

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/queue-set.yaml
```

Look for:

- `spec.queueName` on every item;
- `spec.throughput.maxReadOps` and `maxWriteOps` when provided; and
- `spec.retention.hours`, including its default.

Source: [`components/queue_set.go`](../components/queue_set.go)

## Route catalog

**Explore:** `ForEachMap.WithBody` producing a structured object under every
source key.

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/route-catalog.yaml
```

Expected shape:

```yaml
spec:
  routes:
    checkout:
      service: checkout-api
      metadata:
        sourceRoute: checkout
    search:
      service: search-api
      metadata:
        sourceRoute: search
```

Source: [`components/route_catalog.go`](../components/route_catalog.go)

## Managed record

**Explore:** a health policy that treats missing status as not ready.

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/managed-record.yaml
```

Offline dry-run renders the resource but does not simulate status transitions.
Inspect the generated health policy instead:

```bash
sed -n '/healthPolicy/,/"""#/p' generated/component/managed-record.cue
```

Look for condition sources that default to empty, along with checks for both
`Ready` and `Synced`.

Source: [`components/managed_record.go`](../components/managed_record.go)

## File share

**Explore:** a typed map-of-objects parameter followed by a guarded map-to-list
transformation.

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

Expected shape:

```yaml
spec:
  mountPoints:
  - name: cache
    path: /cache
    permissions: "0755"
  - name: reports
    path: /reports
    permissions: "0750"
```

Map iteration order is not guaranteed. `cache` and `reports` may appear in the
opposite order.

Source: [`components/file_share.go`](../components/file_share.go)

## Zone replicas

**Explore:** a validator message that includes the rejected field value.

First render a valid Application:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/zone-replica.yaml
```

Then run the intentionally invalid Application:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml
```

The second command exits non-zero. Confirm that the output contains:

```text
zone 'us-east-1' must not match the primary zone
```

The surrounding CUE error text can vary by CLI version; the value-specific
message is the part expected to remain consistent.

Source: [`components/zone_replica.go`](../components/zone_replica.go)

## Composite store

**Explore:** a primary output, a named output, a repeated output group, and one
health policy covering all of them.

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/composite-store.yaml
```

The command renders:

- `Store/orders`;
- `AccessPolicy/orders-policy`;
- `AccessPoint/orders-ap-1`; and
- `AccessPoint/orders-ap-2`.

Inspect the typed health policy in
[`components/composite_store.go`](../components/composite_store.go). It checks the
primary store, the named policy, and every output with the `accessPoint` prefix.

## Tenant space

**Explore:** `NotMatches` in both validation and conditional rendering.

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/tenant-space.yaml
```

The example tenant does not start with `cust-`, so the output contains:

```yaml
spec:
  tenantName: platform-services
  tier: internal
```

Source: [`components/tenant_space.go`](../components/tenant_space.go)

## Validate all examples together

```bash
./scripts/verify.sh
```

This command regenerates the definitions, compiles the Go packages, validates
every CUE definition, renders every valid Application, and confirms that the
intentionally invalid Application is rejected.
