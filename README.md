# KubeVela DefKit Enhancement Playground

Explore eight DefKit capabilities added to KubeVela through runnable
component-authoring examples. No cluster is required.

With these examples, you can:

- read Go-authored `ComponentDefinition` examples;
- generate the CUE definitions produced by DefKit;
- validate the generated CUE;
- render example Applications with `vela dry-run --offline`; and
- compare each DefKit API with the resource shape it produces.

The examples use illustrative `example.com/v1alpha1` resources, so you do not
need Kubernetes, CRDs, or controllers.

## Quick start

### Requirements

- Go 1.23.8 or newer
- KubeVela CLI (`vela`) with `def vet` and offline `dry-run` support
- Git

### Clone and verify

```bash
git clone https://github.com/naveenkumarsangi/kubevela-community-playground.git
cd kubevela-community-playground

go mod download
./scripts/verify.sh
```

A successful run ends with:

```text
rejected with: zone 'us-east-1' must not match the primary zone
Playground verification passed.
```

This rejection is expected. It confirms that the generated validator reports the
specific value that violated the rule.

## Try your first example

First, generate all ComponentDefinitions:

```bash
go run ./cmd/generate generated
```

Render the `file-share` example:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

Look for:

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

The Application supplies mount points as a map keyed by name. The DefKit template
converts the map to a list, copies each key into `name`, and adds the default
permission only when the input omits it.

Next, follow the [getting-started guide](docs/getting-started.md), then use the
[example catalog](docs/examples.md) to try every capability.

## Included capabilities

| What you can explore | DefKit API or behavior | Example | Related KubeVela issue and pull request |
| --- | --- | --- | --- |
| Write nested fields inside generated list items | Dotted paths in `ItemBuilder.Set` | `queue-set` | [#7282](https://github.com/kubevela/kubevela/issues/7282) / [#7302](https://github.com/kubevela/kubevela/pull/7302) |
| Build structured values for every map entry | `ForEachMap.WithBody` | `route-catalog` | [#7283](https://github.com/kubevela/kubevela/issues/7283) / [#7339](https://github.com/kubevela/kubevela/pull/7339) |
| Treat an absent resource status as not ready | Status-safe health expressions | `managed-record` | [#7284](https://github.com/kubevela/kubevela/issues/7284) / [#7299](https://github.com/kubevela/kubevela/pull/7299) |
| Define objects under dynamic map keys with typed fields | `Map.OfObject` and `Map.OfSchemaRef` | `file-share` | [#7287](https://github.com/kubevela/kubevela/issues/7287) / [#7314](https://github.com/kubevela/kubevela/pull/7314) |
| Turn a map into a list while retaining key and value | `ForEachMapWith` and `ForEachMapWithGuarded` | `file-share` | [#7288](https://github.com/kubevela/kubevela/issues/7288) / [#7330](https://github.com/kubevela/kubevela/pull/7330) |
| Include the rejected value in a validator message | `ValidateValue` with interpolation | `zone-replica` | [#7289](https://github.com/kubevela/kubevela/issues/7289) / [#7348](https://github.com/kubevela/kubevela/pull/7348) |
| Check health across primary, named, and grouped outputs | `Health.At` and `Health.Every` | `composite-store` | [#7290](https://github.com/kubevela/kubevela/issues/7290) / [#7332](https://github.com/kubevela/kubevela/pull/7332) |
| Express negative regex conditions directly | `NotMatches` | `tenant-space` | [#7353](https://github.com/kubevela/kubevela/issues/7353) / [#7352](https://github.com/kubevela/kubevela/pull/7352) |

For each capability, [DefKit enhancements](docs/enhancements.md) explains the
motivation, API, command, and expected result.

## How the playground works

```text
Go ComponentDefinition
        │
        │  go run ./cmd/generate
        ▼
Generated CUE definition
        │
        ├── vela def vet
        │
        └── vela dry-run --offline + example Application
                    │
                    ▼
              Rendered resource
```

Each example component calls `defkit.Register`. The generator imports the
component package, reads the resulting registry, and writes one CUE file for each
definition in `generated/component/`.

The verification script runs this full workflow for every example.

## Documentation

- [Getting started](docs/getting-started.md) — prerequisites, setup, generation,
  validation, and rendering
- [Enhancements](docs/enhancements.md) — all eight capabilities with API examples
  and upstream links
- [Example catalog](docs/examples.md) — commands and expected output for every
  example
- [Troubleshooting](docs/troubleshooting.md) — common setup and rendering issues
- [Official DefKit overview](https://kubevela.io/docs/platform-engineers/defkit/overview/)

## Repository layout

```text
components/          Go-authored ComponentDefinitions
examples/            Applications used for offline rendering
cmd/generate/        CUE generator for registered definitions
cmd/register/        JSON registry output
scripts/verify.sh    end-to-end local verification
docs/                user guides and reference material
module.yaml          DefinitionModule metadata
```

## Version note

The module pins this KubeVela revision:

```text
v1.11.1-0.20260917092533-97d67561f080
```

That revision contains all eight demonstrated capabilities. It is a post-merge
pseudo-version rather than a tagged KubeVela release. For production adoption,
use the first tagged release that contains the required changes or pin a revision
you have verified in your own build pipeline.

## Scope

The example resource kinds are illustrative and do not have controllers. Use this
repository to author, validate, and render them offline; do not apply them to a
cluster.

## License

Apache License 2.0. See [LICENSE](LICENSE).
