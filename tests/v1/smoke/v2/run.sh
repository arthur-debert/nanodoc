#!/bin/bash
set -e

# Create output directory
OUTPUT_DIR="$(pwd)/tests/smoke/v2/output"
mkdir -p "$OUTPUT_DIR"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FIXTURES_DIR="$(pwd)/tests/fixtures"

echo "Running smoke tests for nanodoc v2..."
echo "Output will be saved to $OUTPUT_DIR"

# Basic test - process a single Python file
echo "Test 1: Basic file processing"
python -m nanodoc --use-v2 src/nanodoc/core.py > "$OUTPUT_DIR/basic.txt"

# Process with line numbers
echo "Test 2: Line numbers"
python -m nanodoc --use-v2 -n src/nanodoc/core.py > "$OUTPUT_DIR/line_numbers.txt"

# Process multiple files
echo "Test 3: Multiple files"
python -m nanodoc --use-v2 src/nanodoc/core.py src/nanodoc/cli.py > "$OUTPUT_DIR/multiple_files.txt"

# Process with multiple files in directory
echo "Test 4: Multiple files in directory"
python -m nanodoc --use-v2 src/nanodoc/core.py src/nanodoc/cli.py src/nanodoc/toc.py > "$OUTPUT_DIR/multiple_directory_files.txt"

# Process a directory
echo "Test 5: Directory"
python -m nanodoc --use-v2 src/nanodoc/themes > "$OUTPUT_DIR/directory.txt"

# Test theming (v2 specific feature)
echo "Test 6: Theming"
python -m nanodoc --use-v2 --theme neutral src/nanodoc/core.py > "$OUTPUT_DIR/theme_neutral.txt"

# Test TOC generation (v2 feature)
echo "Test 7: TOC generation"
python -m nanodoc --use-v2 --toc src/nanodoc/core.py > "$OUTPUT_DIR/toc.txt"

# Test bundling with sample fixtures
echo "Test 8: Bundle file"
# Create a test directory
BUNDLE_DIR="$OUTPUT_DIR/bundle"
mkdir -p "$BUNDLE_DIR"

# Create a main bundle file
cat > "$BUNDLE_DIR/main.ndoc" << EOF
# Main file
def main():
    print("Hello world")

# @include sub.ndoc
EOF

# Create the included file
cat > "$BUNDLE_DIR/sub.ndoc" << EOF
# Sub module
def helper():
    return "I'm a helper function"
EOF

# Process the bundle
echo "Processing bundle file..."
python -m nanodoc --use-v2 "$BUNDLE_DIR/main.ndoc" > "$OUTPUT_DIR/bundle.txt"

# Process with both TOC and theme
echo "Processing bundle with TOC and theme..."
python -m nanodoc --use-v2 --toc --theme neutral "$BUNDLE_DIR/main.ndoc" > "$OUTPUT_DIR/bundle_toc_theme.txt"

echo "All tests completed. Results are in $OUTPUT_DIR"
