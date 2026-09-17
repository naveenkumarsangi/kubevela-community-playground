# Troubleshooting

## `vela: command not found`

Install the KubeVela CLI from the
[KubeVela releases](https://github.com/kubevela/kubevela/releases), or point the
repository scripts at an existing binary:

```bash
export VELA_BIN=/absolute/path/to/vela
./scripts/verify.sh
```

Commands shown individually in the documentation use `vela`. Replace it with the
absolute path or `${VELA_BIN}` when needed.

## `go mod download` cannot resolve KubeVela

Confirm that your machine can reach the configured Go proxy:

```bash
go env GOPROXY
```

Then retry:

```bash
go clean -modcache
go mod download
```

`go clean -modcache` removes the shared module cache and forces every dependency
to download again. Use it only when the cached module is actually corrupt; it is
not a routine setup step.

## The generated directory is missing

Generate the definitions before running `vela def vet` or `vela dry-run`:

```bash
go run ./cmd/generate generated
```

The `generated/` directory is intentionally ignored by Git.

## `vela def vet` mentions external CUE packages

Some CLI versions try to discover external CUE packages before validating a
local definition. Check the final line to see whether validation succeeded:

```text
Validation generated/component/<name>.cue succeed.
```

If the final line does not report successful validation, rerun the full verifier
and keep the complete output:

```bash
./scripts/verify.sh
```

## The invalid replica example returns an error

That is expected:

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/zone-replica-invalid.yaml
```

This example demonstrates a validator message that identifies the rejected
value. Confirm that the output includes:

```text
zone 'us-east-1' must not match the primary zone
```

The surrounding CUE error and the order of `true` and `false` can vary by CLI
version.

## Map entries render in a different order

Map iteration order does not affect the expected results. `file-share` and
`route-catalog` may render their entries in a different order while producing
the same data.

When writing assertions or downstream logic, compare map contents without
assuming a specific iteration order.

## Dry-run cannot find a component type

Confirm that you generated the definitions and passed the component directory:

```bash
go run ./cmd/generate generated

vela dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

Also confirm that the expected CUE file exists under `generated/component/`.

## A command tries to contact a cluster

Every rendering command in this repository uses:

```text
--offline
```

Keep that flag when running `vela dry-run`. The example resources have no CRDs or
controllers and should not be applied to a cluster.

## The pinned revision is unexpected

Check the selected module version:

```bash
go list -m -f 'Module: {{.Path}}{{"\n"}}Version: {{.Version}}' \
  github.com/oam-dev/kubevela
```

Expected:

```text
Module: github.com/oam-dev/kubevela
Version: v1.11.1-0.20260917092533-97d67561f080
```

If it differs, restore `go.mod` and `go.sum` from the repository before running
the examples.

## Verification changed `go.mod` or `go.sum`

The verifier does not modify dependencies. Check the working tree:

```bash
git status --short
```

If the module files changed, inspect the diff before restoring them. Avoid
running `go get` or `go mod tidy` when you only want to try the pinned examples.

## Still stuck?

Capture these commands and their complete output when opening a repository issue:

```bash
go version
vela version
go list -m github.com/oam-dev/kubevela
./scripts/verify.sh
```
