# Changelog

## Unreleased

- Added `MustNewStatus` for package-level custom status definitions that must validate code format.
- Changed `Wrap`, `Wrapf`, and `WithStack` to use `StatusInternalServer()` when wrapping ordinary errors without an existing `goerr.Item`, avoiding accidental success codes in unified response layers.
- Clarified README guidance around `New`, `WrapStatus`, `Wrap`, custom statuses, and business-specific HTTP status mapping.
