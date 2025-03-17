# Nanodoc Redesign Rationale

## Current State and Problems

Nanodoc, at its core, is a file bundler with added features. The current
codebase has grown organically, leading to:

- **Blurred Stage Separation:** The stages of processing (resolving paths,
  resolving files, gathering content, building content, and rendering) are not
  clearly defined, making the code harder to understand and modify.
- **Intermingled Responsibilities:** File resolving is mixed with content
  gathering, and formatting/help are intertwined with core logic.
- **Lack of Clear Data Structures:** The absence of well-defined data structures
  makes it difficult to manage the state of the bundling process and to
  implement new features. Specifically, there's no good way to represent content
  selections from files (ranges) and to track the origin of content (especially
  from inline bundles).

## Proposed Design

The redesigned Nanodoc will follow a clear, staged architecture:

1. **Resolving Paths:** Handles user input (globs, directories, file paths) and
   produces a list of absolute file paths.
2. **Resolving Files:** Takes the file paths and creates initial `FileContent`
   objects. This stage determines if a file is a bundle or a regular file.
3. **Gathering Content:** Reads the actual file contents and applies any line
   range selections, populating the `FileContent` objects.
4. **Building the Content:** Processes bundle files (including handling inline
   directives recursively), flattens the structure into a single list of
   `FileContent` objects, and prepares the data for rendering. This stage
   creates the final `Document` object.
5. **Rendering:** Takes the processed `Document` object and generates the final
   output, handling formatting, TOC generation, and theming.

## Data Structures

The following data structures are key to the new design:

- **`Range`:** Represents a line range within a file.

  ```python
  class Range:
      start: int  # Inclusive start line number
      end: int | None  # Inclusive end line number (None for EOF)
  ```

  Alternatively, a simple tuple `(start, end)` can be used.

- **`FileContent`:** Encapsulates the content and metadata for a single file.

  ```python
  class FileContent:
      filepath: str  # Path to the file
      ranges: list[Range]  # Line ranges to include
      content: str  # Content after applying ranges
      is_bundle: bool  # True if this represents a bundle file
      original_source: str | None  # Source file if part of inline bundle
  ```

- **`Document`:** Represents the entire document after processing bundles.

  ```python
  class Document:
      content_items: list[FileContent]  # Ordered list of content blocks
      toc: list  # TOC data (structure to be defined later)
  ```

  The `toc` could be a list of `(filepath, heading, line_number)` tuples.

## Stage Details

**1. Resolving Paths:**

- **Input:** User arguments (command-line input).
- **Output:** List of absolute file paths (strings).
- **Responsibilities:**
  - Expand globs (e.g., `*.txt`).
  - Handle directory inputs (recursively find files).
  - Validate that paths exist.
  - Convert relative paths to absolute paths.

**2. Resolving Files:**

- **Input:** List of absolute file paths.
- **Output:** List of _initial_ `FileContent` objects.
- **Responsibilities:**
  - Create a `FileContent` object for each path.
  - Determine if the file is a bundle (`.ndoc` extension, or configurable).
  - Initialize `ranges`:
    - If the path includes a range specifier (e.g., `file.txt:10-20`), parse it.
    - Otherwise, default to the entire file (`Range(1, None)`).
  - At this stage, `content` can be empty or contain the _entire_ file's content
    (for efficiency, deferring loading to the next stage is preferable).
  - Set `is_bundle` appropriately.
  - `original_source` is initially `None`.

**3. Gathering Content:**

- **Input:** List of initial `FileContent` objects.
- **Output:** List of `FileContent` objects with the `content` field populated.
- **Responsibilities:**
  - Read the content of each file.
  - Apply the specified `ranges` to extract the relevant lines.
  - Store the extracted content in the `content` field of the `FileContent`
    object.

**4. Building the Content:**

- **Input:** List of `FileContent` objects (with content loaded).
- **Output:** A single `Document` object.
- **Responsibilities:**
  - Process bundle files:
    - Parse bundle directives (e.g., `inline`, `include`).
    - Recursively process included files or bundles (handling potential circular
      dependencies).
    - Create new `FileContent` objects for inlined content, setting
      `original_source` to the bundle's path.
  - Flatten the structure: Combine the `FileContent` objects from the main input
    and any processed bundles into a single, ordered list in the `Document`'s
    `content_items` field.
  - The `toc` field of the `Document` can be populated here (or during
    rendering).

**5. Rendering:**

- **Input:** `Document` object.
- **Output:** Final rendered output (string).
- **Responsibilities:**
  - Generate the Table of Contents (if requested), using the `content_items` and
    potentially extracting headings from the `content`.
  - Apply theming and styling.
  - Add line numbers (if requested).
  - Concatenate the `content` from each `FileContent` in `content_items`,
    handling differences based on `original_source` (e.g., no file separators
    for inlined content).
  - Handle any other formatting options.

## Advantages of this Design

- **Clear Separation of Concerns:** Each stage has a well-defined purpose and
  operates on specific data structures.
- **Improved Maintainability:** Changes in one stage are less likely to affect
  other stages.
- **Extensibility:** New features (e.g., new output formats, different bundle
  directives) can be added more easily.
- **Testability:** Each stage can be tested independently.
- **Handles Inline Bundles Correctly:** The `original_source` field preserves
  context for proper rendering.

## Risks and Considerations

- **Recursive Bundles:** Implement robust circular dependency detection.
- **Error Handling:** Define clear error handling for file I/O, invalid ranges,
  and parsing errors.
- **Performance:** Optimize file reading and content processing for large files
  and complex bundles.
- **Bundle Directive Syntax:** Formalize the syntax for bundle directives.
- **Overlapping Ranges:** Decide how to handle overlapping or conflicting ranges
  within a single file.
