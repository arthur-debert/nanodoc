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
- **Completion Requirements:**
  - Fix path resolution in CLI context
  - Ensure path resolver properly handles command-line arguments

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
- **Completion Requirements:**
  - Verify file resolving works properly with CLI-provided paths
  - Test with various range specifiers through the CLI

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
- **Completion Requirements:**
  - Add end-to-end tests for bundle processing through CLI
  - Test directive parsing with real files through CLI commands

## Milestone 4: Rendering and TOC Generation

- **Deliverables:**
  - Implement the "Rendering" stage:
    - Function to take a `Document` object and generate the final output string.
    - Basic concatenation of `FileContent` content.
    - Implementation of TOC generation (define the `toc` data structure).
    - Unit tests for rendering, including tests for TOC generation.
- **Focus:** Produce the final output from the processed document structure.
- **Completion Requirements:**
  - Fix TOC generation when invoked through CLI
  - Add end-to-end tests for TOC generation
  - Verify TOC appears correctly in output

## Milestone 5: Formatting, Theming, and Options

- **Deliverables:**
  - Implement formatting options (line numbers, etc.).
  - Implement theming/styling capabilities.
  - Refactor existing code to integrate with the new architecture.
  - Comprehensive integration tests.
- **Focus:** Add the finishing touches and ensure the redesigned Nanodoc meets
  all functional requirements.
- **Completion Requirements:**
  - Fix line number formatting in CLI context
  - Ensure theming works properly through CLI
  - Add tests for all formatting options through CLI

## Milestone 6: CLI Integration

- **Deliverables:**
  - Implement CLI integration
  - Test CLI integration
- **Focus:** Make the new implementation accessible through the command line.
- **Completion Requirements:**
  - Fix argument handling in CLI
  - Ensure all options are properly processed
  - Add end-to-end tests for CLI functionality
  - Test each feature (TOC, line numbers, theming) through CLI

## Milestone 7: Documentation and Cleanup

- **Deliverables:**
  - Refactor existing code to integrate with the new architecture.
  - Comprehensive integration tests.
  - Update documentation (README, usage instructions).
  - Code cleanup and refactoring.
  - Final testing and bug fixes.
- **Focus:** Prepare the redesigned Nanodoc for release.
- **Completion Requirements:**
  - Comprehensive documentation for all features
  - Complete test coverage including end-to-end tests
  - Code cleanup and final refactoring
  - Decision on whether v2 should become the default implementation

## Progress Tracking

| Milestone | Description                            | Status      | Completion Date |
| --------- | -------------------------------------- | ----------- | --------------- |
| 1         | Core Data Structures & Path Resolution | Completed   | 2024-03-17      |
| 2         | File Resolving & Content Gathering     | Completed   | 2024-03-17      |
| 3         | Building the Document                  | Completed   | 2024-05-02      |
| 4         | Rendering and TOC Generation           | Completed   | 2024-05-05      |
| 5         | Formatting, Theming, and Options       | Completed   | 2024-05-05      |
| 6         | CLI Integration                        | Completed   | 2024-05-09      |
| 7         | Documentation and Cleanup              | Not Started |                 |

## Testing Strategy

Each stage will include comprehensive tests:

1. **Unit Tests**: Test individual functions in isolation
2. **Integration Tests**: Test interactions between modules
3. **End-to-End Tests**: Test complete pipeline with various inputs
4. **Edge Cases**: Test boundary conditions and error handling

## Reality Check: What's Actually Working?

| Milestone | Component                | Working Status | Issues Identified          |
| --------- | ------------------------ | -------------- | -------------------------- |
| 1         | Core Data Structures     | ✅ Working     | None significant           |
| 1         | Path Resolution          | ✅ Working     | Fixed argument handling    |
| 2         | File Resolving           | ✅ Working     | None significant           |
| 2         | Content Gathering        | ✅ Working     | None significant           |
| 2         | Range Parsing            | ✅ Working     | None significant           |
| 3         | Bundle Processing        | ✅ Working     | Tested with CLI            |
| 3         | Directive Parsing        | ✅ Working     | Tested with CLI            |
| 3         | Document Building        | ✅ Working     | None significant           |
| 4         | Basic Rendering          | ✅ Working     | Works in all contexts      |
| 4         | TOC Generation           | ✅ Working     | Fixed and working with CLI |
| 5         | Line Numbers             | ✅ Working     | Works in all contexts      |
| 5         | Theming                  | ✅ Working     | Works in all contexts      |
| 6         | CLI Integration          | ✅ Working     | Fixed argument handling    |
| 6         | CLI Option Processing    | ✅ Working     | All options work correctly |
| 7         | Documentation            | ❌ Not Started | Planned for future         |
| 7         | End-to-End Tests         | ✅ Added       | Comprehensive tests added  |
| 7         | File Header Format       | ✅ Fixed       | Now matches v1 format      |
| 7         | Bundle File Recognition  | ✅ Enhanced    | Now detects all variants   |
| 7         | v1-v2 Output Consistency | ✅ Verified    | Smoke tests confirm match  |
