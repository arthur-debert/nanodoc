# Nanodoc Redesign Implementation Milestones

This document outlines the implementation milestones for the Nanodoc redesign.
Each milestone has clear deliverables and focuses on a specific part of the new
architecture.

## Milestone 1: Core Data Structures and Path Resolution

- **Deliverables:**
  - Implement the `Range` (or tuple) and `FileContent` data structures.
  - Implement the "Resolving Paths" stage:
    - Function to take user input (globs, paths, directories) and return a list
      of absolute file paths.
    - Unit tests for path resolution, covering various cases (valid paths,
      invalid paths, globs, directories).
- **Focus:** Establish the foundation for data representation and initial file
  processing.

## Milestone 2: File Resolving and Content Gathering

- **Deliverables:**
  - Implement the "Resolving Files" stage:
    - Function to take a list of file paths and create initial `FileContent`
      objects.
    - Logic to determine if a file is a bundle.
    - Initial range parsing (from file path specifiers).
  - Implement the "Gathering Content" stage:
    - Function to read file contents and apply line ranges, populating the
      `content` field of `FileContent` objects.
    - Unit tests for file resolving and content gathering, including tests for
      different range selections.
- **Focus:** Complete the initial processing of files and the extraction of
  their content.

## Milestone 3: Building the Document (Bundle Handling)

- **Deliverables:**
  - Implement the "Building the Content" stage:
    - Function to process bundle files (recursively).
    - Parsing of `inline` and `include` directives (define the syntax clearly).
    - Creation of new `FileContent` objects for inlined content.
    - Flattening the structure into a `Document` object.
    - Circular dependency detection.
    - Unit tests for bundle processing, including tests for nested bundles and
      circular dependencies.
- **Focus:** Handle the core logic of combining content from multiple files and
  bundles.

## Milestone 4: Rendering and TOC Generation

- **Deliverables:**
  - Implement the "Rendering" stage:
    - Function to take a `Document` object and generate the final output string.
    - Basic concatenation of `FileContent` content.
    - Implementation of TOC generation (define the `toc` data structure).
    - Unit tests for rendering, including tests for TOC generation.
- **Focus:** Produce the final output from the processed document structure.

## Milestone 5: Formatting, Theming, and Options

- **Deliverables:**
  - Implement formatting options (line numbers, etc.).
  - Implement theming/styling capabilities.
  - Refactor existing code to integrate with the new architecture.
  - Comprehensive integration tests.
- **Focus:** Add the finishing touches and ensure the redesigned Nanodoc meets
  all functional requirements.

## Milestone 6: Documentation and Cleanup

- **Deliverables:**
  - Update documentation (README, usage instructions).
  - Code cleanup and refactoring.
  - Final testing and bug fixes.
- **Focus:** Prepare the redesigned Nanodoc for release.

## Progress Tracking

| Stage | Description                    | Status      | Completion Date |
| ----- | ------------------------------ | ----------- | --------------- |
| 0     | Setup and Infrastructure       | Not Started |                 |
| 1     | Path Resolution                | Not Started |                 |
| 2     | Content Extraction             | Not Started |                 |
| 3     | Document Construction          | Not Started |                 |
| 4     | Content Formatting             | Not Started |                 |
| 5     | Rendering                      | Not Started |                 |
| 6     | Pipeline Integration           | Not Started |                 |
| 7     | CLI Integration                | Not Started |                 |
| 8     | Documentation and Finalization | Not Started |                 |

## Testing Strategy

Each stage will include comprehensive tests:

1. **Unit Tests**: Test individual functions in isolation
2. **Integration Tests**: Test interactions between modules
3. **End-to-End Tests**: Test complete pipeline with various inputs
4. **Edge Cases**: Test boundary conditions and error handling

## Risk Management

| Risk                                | Mitigation                                                                    |
| ----------------------------------- | ----------------------------------------------------------------------------- |
| Complex document tree construction  | Start with simple cases, then add complexity incrementally                    |
| Preserving context through pipeline | Use immutable data structures and ensure context is passed through each stage |
| Performance with large files        | Implement early performance testing with large inputs                         |
| Backward compatibility              | Clearly document changes and provide migration path                           |
