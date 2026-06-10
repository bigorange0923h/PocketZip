# PocketUnzip password extraction fix

Date: 2026-06-11

## Problem

PocketUnzip could fail to detect password-protected archives correctly.

The archive layer streamed 7z stdout/stderr to the frontend log panel, but it did not keep that output in the returned result. The app layer only stored `ExitErr.Error()`, which is usually a generic value such as `exit status 2`. As a result, `archive.IsPasswordError` often had no real 7z message to inspect, so the app could return a generic extraction failure instead of `password required`.

The app also exposed `GetPasswordCandidates`, but the main `Extract` path did not automatically try saved password candidates before asking the user to enter a password.

## Root Cause

- `archive.Extract` discarded stdout/stderr after emitting log lines.
- Password detection was run against the generic process exit error instead of the captured 7z output.
- Saved password candidates were only surfaced to the frontend dialog, not used by the backend extraction flow.
- `archive.Extract` also assumed a non-nil context, which made direct tests or non-Wails callers panic.

## Fix

- Added `Output` to `archive.ExtractResult`.
- Added `ExtractResult.Error()` so callers can preserve both the process exit status and the 7z diagnostic output.
- Updated stdout/stderr scanning to:
  - continue streaming realtime log lines,
  - collect the same lines into the result output,
  - wait for scanner goroutines before returning.
- Expanded password-error detection for common 7z messages, including:
  - `Wrong password`
  - `password is incorrect`
  - `Cannot open encrypted archive`
  - `Can not open encrypted archive`
  - `Data Error in encrypted file`
  - `Enter password`
- Updated `App.Extract` to:
  - try extraction without a password first,
  - detect password-required failures from captured 7z output,
  - load saved password candidates from SQLite,
  - retry candidates without logging the secret,
  - update password success stats when a saved candidate works,
  - record the final result in history.
- Updated `App.ExtractWithPassword` to use the same extraction and history-recording helpers.
- Made nil contexts fall back to `context.Background()`.

## Tests

- Added archive unit coverage for common password-related 7z output.
- Added archive unit coverage that `ExtractResult.Error()` includes captured output.
- Added an app-level regression test using the test process as a fake 7z executable:
  - first invocation returns `ERROR: Wrong password`,
  - retry with saved `secret` password returns success,
  - history records a successful extraction that used a password.

## Verification

```powershell
go test ./...
```

Result: passed.
