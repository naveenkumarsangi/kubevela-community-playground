# Getting started

This guide takes you from a fresh clone to a validated, rendered DefKit example.
Everything runs locally. You do not need a Kubernetes cluster.

## 1. Check the prerequisites

You need:

- Go 1.23.8 or newer
- the KubeVela CLI (`vela`) with `def vet` and offline `dry-run` support
- Git

Check the installed versions:

```bash
go version
vela version
```

An empty `Core Version` in `vela version` is fine. The commands in this repository
use offline rendering and do not connect to a KubeVela control plane.

If `vela` is not on `PATH`, set its location before running the scripts:

```bash
export VELA_BIN=/absolute/path/to/vela
```

The verification script honors `VELA_BIN`. For individual commands, replace
`vela` with `${VELA_BIN}` or the absolute path.

## 2. Clone the repository

```bash
git clone https://github.com/naveenkumarsangi/kubevela-community-playground.git
cd kubevela-community-playground
```

## 3. Download Go dependencies

```bash
go mod download
```

The repository pins a KubeVela pseudo-version containing all demonstrated DefKit
APIs. See the [version note](../README.md#version-note) before reusing the pin in
another project.

## 4. Run the complete verification

```bash
./scripts/verify.sh
```

The script performs five checks:

1. generates all seven ComponentDefinitions (`file-share` demonstrates two capabilities);
2. compiles every Go package;
3. validates every generated CUE definition with `vela def vet`;
4. renders every valid Application with `vela dry-run --offline`; and
5. confirms that the intentionally invalid replica example includes the rejected
   zone in its message.

A successful run ends with:

```text
rejected with: zone 'us-east-1' must not match the primary zone
Playground verification passed.
```

## 5. Generate the definitions yourself

```bash
go run ./cmd/generate generated
```

The command writes:

```text
generated/component/queue-set.cue
generated/component/route-catalog.cue
generated/component/managed-record.cue
generated/component/file-share.cue
generated/component/zone-replica.cue
generated/component/composite-store.cue
generated/component/tenant-space.cue
```

The registry may print these files in a different order. This is expected; the
order is not guaranteed.

## 6. Validate one generated definition

```bash
vela def vet generated/component/file-share.cue
```

Expected result:

```text
Validation generated/component/file-share.cue succeed.
```

Some CLI versions may first log a warning while looking for external CUE packages.
For these examples, validation has succeeded if the final line is
`Validation generated/component/file-share.cue succeed.`

## 7. Render one Application

```bash
vela dry-run --offline \
  -d generated/component \
  -f examples/file-share.yaml
```

Inspect `spec.mountPoints` in the rendered `FileShare` resource. The map keys from
the Application become list-item names, and the omitted permission defaults to
`"0755"`.

## 8. Explore the rest

- Use the [example catalog](examples.md) for commands and expected output.
- Read [DefKit enhancements](enhancements.md) to understand the APIs behind each
  example.
- Check [troubleshooting](troubleshooting.md) if a command behaves differently on
  your machine.

## Reset generated output

Generated definitions are ignored by Git. Remove and recreate them at any time:

```bash
rm -rf generated
go run ./cmd/generate generated
```
