# GPT-Load (Fork)

Upstream: https://github.com/tbphp/gpt-load

This fork keeps only the differences from upstream.

## Differences

- UI overhaul (Apple-like minimal visuals, typography, spacing, and shadows).
- Login page refresh with a text-only logo treatment.
- New Model Usage page with donut chart, hover tooltip, and free date-range filter.
- Model usage stats aggregated from request_logs with success/error/retry counts; empty models filtered.
- Key management UX polish (view-key modal, layout/scroll fixes, button/input styling).
- Removed unified OpenAI aggregate proxy; use per-channel/per-group proxy endpoints.
