# Nanodoc Smoke Tests

This directory contains smoke tests for the Nanodoc CLI. These tests verify that
basic functionality works correctly in both v1 and v2 implementations.

## Running the Tests

To run the smoke tests:

```bash
python tests/smoke/test_smoke.py
```

## Tests Included

The smoke tests verify:

1. Basic file concatenation
2. Table of Contents generation
3. Line numbering
4. Bundle file processing
5. Theme application

## Adding New Tests

To add a new smoke test:

1. Add a new function to `test_smoke.py` with the test logic
2. Add the function to the `tests` list in the `main()` function

## Troubleshooting

If a test fails, the output will show:

- The exact command that was run
- The stdout and stderr of the command
- The specific assertion that failed

The smoke tests include proper error handling and cleanup of temporary files.
