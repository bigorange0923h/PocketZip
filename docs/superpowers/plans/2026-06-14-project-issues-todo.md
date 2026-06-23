# PocketZip project issues TODO

Date: 2026-06-14

## Current Baseline

- `npm run build` passes.
- `go test ./...` passes after building `frontend/dist` first.
- A small reliability fix was applied during this review:
  - `internal/archive/compress.go` no longer uses dynamic strings as `fmt.Errorf` format strings.
  - `history.ExtractHistory` and `password.PasswordRecord` now expose camelCase JSON fields that match the frontend.
  - `password.PasswordRecord.EncryptedPassword` is no longer exposed through Wails models.

## P0 - Password Save And Reuse Flow

### Problem

The UI exposes a "remember password" concept, but the main extraction flow does not consistently save manually entered passwords after successful extraction. `InlinePasswordPanel.vue` also expects `record.password`, but the backend should not return plaintext passwords.

### Fix Plan

- Remove plaintext password display and selection from `InlinePasswordPanel.vue`.
- Keep saved passwords backend-only for automatic retry.
- Add a manual password submit flow that passes both:
  - password
  - rememberPassword
- Save the password only after a successful manual extraction and only when `rememberPassword` is true.
- Add password deduplication in `password.Save`:
  - match by archive path or archive name,
  - decrypt existing records,
  - update usage metadata instead of inserting duplicate records when the password already exists.
- Update password manager UI to show only metadata:
  - archive name,
  - archive path,
  - success count,
  - last used time,
  - created time.

### Acceptance Criteria

- Manual password extraction succeeds without exposing saved plaintext passwords.
- When "remember password" is checked, successful password is saved once.
- Reusing the same password does not create duplicate rows.
- Password manager does not contain plaintext or encrypted password bytes in frontend models.

## P0 - Strategy Persistence And Execution

### Problem

Extraction strategies are currently partial placeholders. `SaveExtractStrategy` stores simplified strings such as `default` or `retry`, so fields like `outputDir`, `maxRetries`, and `autoOpen` are not fully persisted. `autoOpen` is also not executed after successful strategy extraction.

### Fix Plan

- Store each strategy as JSON in `app_config`.
- Validate strategy data:
  - non-empty name,
  - `maxRetries` clamped or rejected outside `0..10`,
  - output directory optional.
- Update `GetExtractStrategy` to parse JSON and return clear errors for invalid stored values.
- Update `GetExtractStrategies` to merge built-in defaults with saved overrides.
- Update `ExtractWithStrategy`:
  - use configured output directory when present,
  - perform retry when `autoRetry` is true,
  - open output directory after success when `autoOpen` is true.
- Add tests for save/load/default/invalid strategy behavior.

### Acceptance Criteria

- Editing strategy settings persists across app restart.
- `autoOpen` actually opens the output directory after successful extraction.
- Invalid strategy JSON does not silently become a default strategy.

## P1 - Real Cancellation

### Problem

The cancel button currently changes frontend state only. The underlying `7z` process continues to run.

### Fix Plan

- Add operation IDs for extract/compress jobs.
- Store cancel functions or process handles in `App`.
- Add Wails methods:
  - `CancelExtract(operationID)`
  - `CancelCompress(operationID)`
- Wire `ProgressOverlay` cancel button to backend cancellation.
- Ensure cancellation records history as a cancelled/failed operation with a clear message.
- Add tests with fake long-running 7z commands.

### Acceptance Criteria

- Clicking cancel terminates the active `7z` process.
- UI stops because backend reports cancellation, not because frontend state was reset.
- No orphaned `7z.exe` remains after cancellation.

## P1 - Archive Path Safety

### Problem

Extraction currently delegates directly to `7z x`. Malicious archives may contain paths that escape the output directory, absolute paths, drive-letter paths, UNC paths, or `..` segments.

### Fix Plan

- Use `7z l -slt` to list archive entries before extraction.
- Parse structured `Path = ...` output instead of column-based `7z l` output.
- Reject entries that are:
  - absolute paths,
  - Windows drive-letter paths,
  - UNC paths,
  - paths containing `..`,
  - paths that normalize outside the output directory.
- Apply validation to:
  - normal extraction,
  - extraction with saved password retry,
  - manual password extraction,
  - batch extraction,
  - strategy extraction.

### Acceptance Criteria

- Unsafe archive paths are rejected before running `7z x`.
- Tests cover safe paths, `..`, drive-letter, UNC, and absolute path cases.

## P1 - UI Text Encoding Cleanup

### Problem

Many Chinese strings in frontend components, docs, and comments are mojibake. This harms user experience and maintenance.

### Fix Plan

- Prioritize user-visible Vue components:
  - `App.vue`
  - `ProgressOverlay.vue`
  - `InlinePasswordPanel.vue`
  - `PasswordDialog.vue`
  - `HistoryList.vue`
  - `PasswordManager.vue`
  - `StrategyManager.vue`
- Replace corrupted strings with readable Chinese.
- Then clean backend comments only where they help maintenance.
- Avoid changing behavior during text-only cleanup.

### Acceptance Criteria

- Main UI has readable Chinese text.
- Frontend build still passes.
- Text-only cleanup is committed separately from behavior changes.

## P1 - Compression Workflow Validation

### Problem

Compression support exists but still needs full functional verification, especially for folders, password-protected archives, `7z` header encryption, existing output files, and `.gz` behavior.

### Fix Plan

- Add backend tests for `buildCompressArgs`.
- Add fake 7z tests for successful compression and failed compression.
- Verify actual Windows behavior manually with:
  - single file to zip,
  - multiple files to zip,
  - folder to 7z,
  - password-protected 7z,
  - existing archive overwrite behavior,
  - `.gz` naming and command semantics.
- Add UI error handling around save-path selection and compression failures.

### Acceptance Criteria

- Compression succeeds for selected files/folders.
- Password-protected 7z archives can be created and later extracted.
- Existing output behavior is explicit and user-friendly.

## P2 - Development Versus Production Options

### Problem

`Debug.OpenInspectorOnStartup` is currently enabled in `main.go`, which is useful during active debugging but undesirable for release builds.

### Fix Plan

- Gate inspector startup behind a build tag or environment/config flag.
- Keep `wails dev` convenient, but ensure production builds do not open DevTools.

### Acceptance Criteria

- Development builds can open DevTools.
- Production builds do not open DevTools by default.

## P2 - Build And Verification Ordering

### Problem

Running `go test ./...` while `npm run build` is replacing hashed files in `frontend/dist` can produce transient embed errors.

### Fix Plan

- Document the verification order:
  1. `npm run build`
  2. `go test ./...`
  3. `wails build`
- Add a local script that runs the sequence in order.
- Avoid running frontend build and Go tests in parallel when the root package embeds `frontend/dist`.

### Acceptance Criteria

- A single script verifies the project reliably.
- Contributors do not hit transient missing hashed asset errors.

## Recommended Implementation Order

1. P0 Password save and reuse flow.
2. P0 Strategy persistence and execution.
3. P1 Real cancellation.
4. P1 Archive path safety.
5. P1 Compression workflow validation.
6. P1 UI text encoding cleanup.
7. P2 Development versus production options.
8. P2 Build and verification ordering.
