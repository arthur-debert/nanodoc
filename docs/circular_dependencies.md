# Circular Dependencies in Nanodoc

This document explains how nanodoc handles circular dependencies in bundle files and live bundles.

## What are Circular Dependencies?

A circular dependency occurs when files reference each other in a loop, creating an infinite cycle. For example:

- File A includes File B
- File B includes File A

This creates an infinite loop that would cause the program to run forever or crash.

## How Nanodoc Detects Circular Dependencies

Nanodoc uses sophisticated tracking to detect circular references:

1. **Bundle Files**: When processing `.bundle.*` files, nanodoc tracks which bundles have been visited
2. **Live Bundles**: When processing `[[file:]]` directives, nanodoc tracks the inclusion chain
3. **Depth Limits**: As a safety measure, live bundles have a maximum nesting depth of 10 levels

## Examples of Circular Dependencies

### Simple Circular Reference

```text
# file1.bundle.txt
file2.bundle.txt

# file2.bundle.txt
file1.bundle.txt
```

### Three-Way Circle

```text
# a.bundle.txt
b.bundle.txt

# b.bundle.txt
c.bundle.txt

# c.bundle.txt
a.bundle.txt
```

### Self-Reference

```text
# recursive.bundle.txt
recursive.bundle.txt
other-file.txt
```

### Live Bundle Circular Reference

```text
# doc1.txt
Content here
[[file:doc2.txt]]

# doc2.txt
More content
[[file:doc1.txt]]
```

## Error Messages

When nanodoc detects a circular dependency, it provides a clear error message showing:

- The file that triggered the detection
- The chain of files that form the cycle

Example error:

```
Error: circular dependency detected: bundle1.txt -> [bundle2.txt, bundle3.txt, bundle1.txt]
```

## Valid Patterns That Are NOT Circular

### Diamond Pattern

Multiple files can include the same file without creating a circular dependency:

```
# top.bundle.txt
left.bundle.txt
right.bundle.txt

# left.bundle.txt
shared.txt

# right.bundle.txt
shared.txt

# shared.txt
Common content
```

This is valid because there's no cycle - just multiple paths to the same file.

### Deep Nesting

Files can be nested deeply as long as they don't reference back:

```
# level1.txt
[[file:level2.txt]]

# level2.txt
[[file:level3.txt]]

# level3.txt
Final content
```

## Best Practices

1. **Plan Your Structure**: Before creating complex bundle hierarchies, plan the structure to avoid cycles
2. **Use Descriptive Names**: Clear file names help identify potential circular references
3. **Keep It Simple**: Deeply nested bundles are harder to maintain and debug
4. **Test Incrementally**: Add files one at a time and test to catch circular dependencies early

## Troubleshooting

If you encounter a circular dependency error:

1. **Check the Error Message**: The error shows the exact chain of files involved
2. **Review Each File**: Open each file in the chain to understand the references
3. **Break the Cycle**: Remove one of the references to break the circular dependency
4. **Consider Restructuring**: Sometimes a circular dependency indicates a need to reorganize your files

## Technical Details

- Bundle files track visited paths using absolute file paths
- Live bundles track visited paths within each processing session
- The maximum depth for live bundles prevents stack overflow from deep nesting
- Error reporting includes the full dependency chain for easy debugging
