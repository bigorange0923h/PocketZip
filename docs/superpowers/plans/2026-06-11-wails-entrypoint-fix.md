# PocketUnzip Wails entrypoint fix

Date: 2026-06-11

## Problem

`wails dev` and the generated executable failed to start with:

```text
This version of %1 is not compatible with the version of Windows you're running.
```

The production build also produced `build/bin/PocketUnzip.exe`, but the file was only about 145 KB and could not be launched.

## Root Cause

Wails builds the package in the project root by default. The project root was `package pocketunzip` and only contained embedded frontend assets. The actual application entrypoint lived in `cmd/pocketunzip/main.go`.

Because the root package was not `package main`, Wails effectively built a Go package archive and wrote it to `PocketUnzip.exe`. The file started with:

```text
!<arch>
```

instead of the Windows executable header:

```text
MZ
```

Windows correctly rejected that file as not being a valid executable for the OS.

## Fix

- Added a root-level `main.go` with the Wails application startup code.
- Changed root `assets.go` from `package pocketunzip` to `package main`, so embedded frontend assets are available to the root Wails app.
- Marked the old `cmd/pocketunzip/main.go` with `//go:build legacycmd` so it no longer participates in normal builds.
- Kept the 7z discovery logic in the new root entrypoint.

## Verification

```powershell
go test ./...
wails build -skipbindings -clean
```

The rebuilt executable is now a real Windows PE file:

```text
Header: MZ
Size:   about 15 MB
```

The app starts successfully in both production build mode and Wails development mode.

For this session, development mode was started with:

```powershell
wails dev -skipbindings
```

`-skipbindings` is acceptable for this change because no exported Wails method signatures changed.
