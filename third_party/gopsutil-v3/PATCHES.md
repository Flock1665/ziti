# gopsutil/v3 — local patches (iOS only)

Vendored copy of `github.com/shirou/gopsutil/v3` v3.24.5 (upstream README.md
and LICENSE kept as-is). Only build constraints were changed, and only to make
`GOOS=ios GOARCH=arm64 CGO_ENABLED=1` compile.

## Why

iOS builds require `CGO_ENABLED=1`, which makes the `darwin && cgo` files in
`cpu/` and `process/` selectable — and iOS satisfies the `darwin` build tag.
Both include headers the iPhoneOS SDK does not ship:

```
cpu/cpu_darwin_cgo.go:15:10:     fatal error: 'libproc.h' file not found
process/process_darwin_cgo.go:7: fatal error: 'libproc.h' file not found
```

The include is guarded by `#if TARGET_OS_MAC`, which is also 1 on iOS — the
desktop-only check would be `TARGET_OS_OSX`. Upstream issue: shirou/gopsutil#1230
(open; v3.24.5 is the latest v3 release).

## Changed files (build constraints only)

| File | Change |
|------|--------|
| `cpu/cpu_darwin_cgo.go` | `darwin && cgo` -> `darwin && cgo && !ios` |
| `cpu/cpu_darwin_nocgo.go` | `darwin && !cgo` -> `darwin && (!cgo \|\| ios)` |
| `process/process_darwin_cgo.go` | `darwin && cgo` -> `darwin && cgo && !ios` |
| `process/process_darwin_nocgo.go` | `darwin && !cgo` -> `darwin && (!cgo \|\| ios)` |

Function-by-function comparison shows the `_nocgo` implementations define the
same entry points as the `_cgo` ones (process-only helpers such as `procArgs`
are private to the cgo file and referenced nowhere else), so the swap is safe.
macOS builds are unchanged and still use the cgo implementations.

`mem/` was audited and left untouched: its darwin cgo file only uses mach APIs
(`mach/mach_host.h`, `mach/vm_page_size.h`), which the iPhoneOS SDK provides.
(the same gopsutil issue lists cpu/disk/docker/host/process as failing on iOS —
mem is not among them; disk/docker/host are not in ziti's build graph.)

## Runtime note

On iOS the `_nocgo` implementations use sysctl-based queries. Under the iOS
sandbox some may be restricted; they are only exercised by the diagnostic
`ziti agent ...` commands, never by `ziti router run`.

## Maintenance

Remove this directory and the `replace` line in the root `go.mod` once
upstream ships an iOS-safe release.
