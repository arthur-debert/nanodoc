#!/bin/bash
set -e

# Create output directory
OUTPUT_DIR="$(pwd)/tests/smoke/v1/output"
mkdir -p "$OUTPUT_DIR"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FIXTURES_DIR="$(pwd)/tests/fixtures"

echo "Running smoke tests for nanodoc v1..."
echo "Output will be saved to $OUTPUT_DIR"

# Basic test - process a single Python file
echo "Test 1: Basic file processing"
python -m nanodoc src/nanodoc/core.py > "$OUTPUT_DIR/basic.txt"

# Process with line numbers
echo "Test 2: Line numbers"
python -m nanodoc -n src/nanodoc/core.py > "$OUTPUT_DIR/line_numbers.txt"

# Process multiple files
echo "Test 3: Multiple files"
python -m nanodoc src/nanodoc/core.py src/nanodoc/cli.py > "$OUTPUT_DIR/multiple_files.txt"

# Process with multiple files instead of glob
echo "Test 4: Multiple files in directory"
python -m nanodoc src/nanodoc/core.py src/nanodoc/cli.py src/nanodoc/toc.py > "$OUTPUT_DIR/multiple_directory_files.txt"

# Process a directory
echo "Test 5: Directory"
python -m nanodoc src/nanodoc/themes > "$OUTPUT_DIR/directory.txt"

# Test bundling with sample fixtures
echo "Test 6: Bundle file"
if [ -d "$FIXTURES_DIR" ]; then
  # If fixtures directory exists, use it
  for fixture in "$FIXTURES_DIR"/*.py; do
    if [ -f "$fixture" ]; then
      base_name=$(basename "$fixture")
      echo "Processing fixture: $base_name"
      python -m nanodoc "$fixture" > "$OUTPUT_DIR/fixture_${base_name}.txt"
    fi
  done
else
  # Otherwise, use our own temporary fixtures
  echo "Creating temporary fixtures..."
  TEMP_FIXTURES="$OUTPUT_DIR/temp_fixtures"
  mkdir -p "$TEMP_FIXTURES"

  # Create a main file that includes another file
  cat > "$TEMP_FIXTURES/main.py" << EOF
# Main file
def main():
    print("Hello world")

# @include sub.py
EOF

  # Create the included file
  cat > "$TEMP_FIXTURES/sub.py" << EOF
# Sub module
def helper():
    return "I'm a helper function"
EOF

  # Process the bundle
  python -m nanodoc "$TEMP_FIXTURES/main.py" > "$OUTPUT_DIR/bundle.txt"
fi

# Create fixtures for testing
echo "Creating fixtures for testing..."
FIXTURES_DIR="$OUTPUT_DIR/fixtures"
mkdir -p "$FIXTURES_DIR"

# Create a main file that includes another file
cat > "$FIXTURES_DIR/main.py" << EOF
# Main file
def main():
    print("Hello world")

# @include sub.py
EOF

# Create the included file
cat > "$FIXTURES_DIR/sub.py" << EOF
# Sub module
def helper():
    return "I'm a helper function"
EOF

# Process the bundle
python -m nanodoc "$FIXTURES_DIR/main.py" > "$OUTPUT_DIR/bundle.txt"

echo "All tests completed. Results are in $OUTPUT_DIR"
