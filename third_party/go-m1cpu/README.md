# go-m1cpu (patched local fork)

Vendored copy of `github.com/shoenig/go-m1cpu` v0.2.2, patched **only** for iOS.

## Why

Building the ziti CLI for `GOOS=ios GOARCH=arm64` requires `CGO_ENABLED=1`
(Go's ios/arm64 target always links externally). With cgo enabled, this
module's `cpu.go` is selected — its cgo preamble does:

```c
#if !defined(MAC_OS_VERSION_12_0) || MAC_OS_X_VERSION_MIN_REQUIRED < MAC_OS_VERSION_12_0
#define kIOMainPortDefault kIOMasterPortDefault
#endif
```

`MAC_OS_VERSION_12_0` is a macOS-only macro, so on the iPhoneOS SDK the guard
is always true and the code ends up calling `kIOMasterPortDefault`, which the
iPhoneOS SDK marks unavailable:

```
error: 'kIOMasterPortDefault' is unavailable: not available on iOS
```

Upstream has no fixed release (v0.2.2 is the latest).

## What changed (2 lines)

| File | Change |
|------|--------|
| `cpu.go` | build constraint `darwin && arm64 && cgo` -> `darwin && arm64 && cgo && !ios` |
| `incompatible.go` | build constraint `!darwin \|\| !arm64 \|\| !cgo` -> `... \|\| ios` |

On iOS the package now resolves to the `incompatible.go` stubs, where
`IsAppleSilicon()` returns false. That is safe at runtime: the only consumer in
ziti's dependency graph is `gopsutil/v3/cpu`, which calls `PCoreHz()` only
inside `if m1cpu.IsAppleSilicon()` — the panicking stub functions are never
reached. The only effect is that CPU MHz is not reported on iOS.

macOS builds are unchanged and still use the real IOKit implementation.

## Maintenance

Remove this directory and the `replace` line in the root `go.mod` once
upstream ships a release that guards the IOKit constant for iOS.

License: MPL-2.0, see `LICENSE` (copied from upstream).
