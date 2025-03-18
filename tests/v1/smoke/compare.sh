#!/bin/bash
set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
V1_DIR="$SCRIPT_DIR/v1"
V2_DIR="$SCRIPT_DIR/v2"
DIFF_DIR="$SCRIPT_DIR/diff"

echo "Running smoke tests and comparing outputs..."

# Create diff directory
mkdir -p "$DIFF_DIR"

# Run v1 tests
echo "Running V1 smoke tests..."
"$V1_DIR/run.sh"

# Run v2 tests
echo -e "\nRunning V2 smoke tests..."
"$V2_DIR/run.sh"

# Compare common output files
echo -e "\nComparing outputs..."
V1_OUTPUT="$V1_DIR/output"
V2_OUTPUT="$V2_DIR/output"

# Get list of common files to compare
COMMON_FILES=""
for v1_file in "$V1_OUTPUT"/*.txt; do
  base_name=$(basename "$v1_file")
  if [ -f "$V2_OUTPUT/$base_name" ]; then
    COMMON_FILES="$COMMON_FILES $base_name"
  fi
done

# Run diffs on common files
for file in $COMMON_FILES; do
  echo -e "\nDiff for $file:"
  diff_output="$DIFF_DIR/${file%.txt}_diff.txt"

  # Use diff with unified format and side-by-side comparison
  if diff -u "$V1_OUTPUT/$file" "$V2_OUTPUT/$file" > "$diff_output"; then
    echo "✅ Files are identical"
  else
    echo "❌ Differences found (saved to $diff_output)"

    # Print a brief summary of changes
    added=$(grep -c "^+" "$diff_output" || true)
    removed=$(grep -c "^-" "$diff_output" || true)
    echo "  - Lines added: $added"
    echo "  - Lines removed: $removed"
  fi
done

# List V2-specific files
echo -e "\nV2-specific output files:"
for v2_file in "$V2_OUTPUT"/*.txt; do
  base_name=$(basename "$v2_file")
  if [ ! -f "$V1_OUTPUT/$base_name" ]; then
    echo "- $base_name"
  fi
done

echo -e "\nDiff analysis completed. Full diff files are in $DIFF_DIR"
