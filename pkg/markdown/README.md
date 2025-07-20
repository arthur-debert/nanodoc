# Markdown Package API Documentation

This package provides markdown parsing, transformation, and rendering capabilities for nanodoc. It supports the requirements for phases 2.1 through 2.4 of markdown support implementation.

## Overview

The markdown package is built on top of:
- **goldmark**: For parsing markdown into AST
- **goldmark-markdown**: For rendering AST back to markdown

## API Components

### 1. Parser
Parses markdown content into an AST structure.

```go
parser := markdown.NewParser()
doc, err := parser.Parse([]byte(markdownContent))
```

### 2. Transformer
Provides methods to modify the markdown AST.

```go
transformer := markdown.NewTransformer()

// Check if document has H1 headers
hasH1 := transformer.HasH1(doc)

// Adjust all header levels (e.g., H1->H2, H2->H3)
err := transformer.AdjustHeaderLevels(doc, 1)

// Insert a file header at the beginning
err := transformer.InsertFileHeader(doc, "filename.md", 2) // Insert as H2
```

### 3. Renderer
Converts the AST back to markdown text using goldmark-markdown.

```go
renderer := markdown.NewRenderer()
output, err := renderer.Render(doc)
```

### 4. TOC Generator
Extracts headers and generates table of contents.

```go
tocGen := markdown.NewTOCGenerator()

// Extract all headers from document
entries := tocGen.ExtractTOC(doc)

// Generate markdown-formatted TOC
tocMarkdown := tocGen.GenerateTOCMarkdown(entries)
```

### 5. Header Formatter
Formats file headers according to nanodoc options.

```go
formatter := markdown.NewHeaderFormatter()
header := formatter.FormatFileHeader("file.md", "1", 2) // "1. file.md"
```

## Usage Examples

### Phase 2.2: Header Level Adjustment
```go
// Parse markdown file
parser := markdown.NewParser()
doc, _ := parser.Parse(content)

// Check if it has H1 and adjust levels if needed
transformer := markdown.NewTransformer()
if transformer.HasH1(doc) {
    transformer.AdjustHeaderLevels(doc, 1)
}

// Render back to markdown
renderer := markdown.NewRenderer()
output, _ := renderer.Render(doc)
```

### Phase 2.3: File Headers Support
```go
// Add file header based on formatting options
headerText := formatter.FormatFileHeader(filename, sequence, 2)
transformer.InsertFileHeader(doc, headerText, 2)
```

### Phase 2.4: Table of Contents
```go
// Extract headers from all documents
var allEntries []markdown.TOCEntry
for _, doc := range documents {
    entries := tocGen.ExtractTOC(doc)
    allEntries = append(allEntries, entries...)
}

// Generate TOC markdown
tocContent := tocGen.GenerateTOCMarkdown(allEntries)
```

## Integration with nanodoc

To integrate this package into nanodoc's renderer:

1. Replace `renderMarkdownBasic()` in `pkg/nanodoc/renderer.go`
2. Use the markdown package to:
   - Parse each markdown file
   - Apply transformations based on FormattingOptions
   - Generate TOC if requested
   - Render final output

Example integration:
```go
func renderMarkdownEnhanced(doc *Document, ctx *FormattingContext) (string, error) {
    parser := markdown.NewParser()
    transformer := markdown.NewTransformer()
    renderer := markdown.NewRenderer()
    
    var processedDocs []*markdown.Document
    
    // Process each content item
    for i, item := range doc.ContentItems {
        // Parse
        mdDoc, err := parser.Parse([]byte(item.Content))
        if err != nil {
            return "", err
        }
        
        // Apply transformations
        if transformer.HasH1(mdDoc) && i > 0 {
            transformer.AdjustHeaderLevels(mdDoc, 1)
        }
        
        // Add file header if needed
        if ctx.ShowFilenames {
            sequence := generateSequence(i+1, doc.FormattingOptions.SequenceStyle)
            headerText := formatter.FormatFileHeader(
                filepath.Base(item.Filepath), 
                sequence, 
                2,
            )
            transformer.InsertFileHeader(mdDoc, headerText, 2)
        }
        
        processedDocs = append(processedDocs, mdDoc)
    }
    
    // Generate TOC if requested
    if ctx.ShowTOC {
        // Extract entries from all docs
        // Generate and prepend TOC
    }
    
    // Render all documents
    var output strings.Builder
    for _, mdDoc := range processedDocs {
        rendered, _ := renderer.Render(mdDoc)
        output.Write(rendered)
        output.WriteString("\n")
    }
    
    return output.String(), nil
}
```

## Testing

The package includes comprehensive tests covering:
- Basic parsing and rendering
- Header level adjustments
- File header insertion
- TOC generation
- Complex markdown documents
- Full integration workflows

Run tests with:
```bash
go test ./pkg/markdown/...
```