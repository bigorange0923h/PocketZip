# PocketZip file selector fix

Date: 2026-06-11

## Problem

Manual validation found two file-selection issues:

- Dragging an archive into the app did not populate the archive path.
- Clicking the file selector did not provide useful feedback when the native file dialog failed to open.

## Root Cause

The old drag-and-drop implementation used browser `DataTransfer.files[0].path`. In WebView2/Wails this value is usually unavailable for security and compatibility reasons, so the frontend could see a dropped file object but not its real filesystem path.

The click handler called the backend `SelectFile` method without surfacing errors in the UI, which made dialog failures look like a dead click.

## Fix

- Replaced browser-only path probing with Wails runtime drag-and-drop support:
  - `OnFileDrop`
  - `OnFileDropOff`
- Marked the selector element as a Wails drop target with `--wails-drop-target: drop`.
- Kept a small browser fallback for environments that still expose `File.path`.
- Added visible error text when file path resolution or native dialog opening fails.
- Restored readable Chinese text in the file selector component.

## Verification

```powershell
npm run build
go test ./...
```

Result: passed.

The Wails development app was restarted after the fix so the running window uses the updated frontend code.
