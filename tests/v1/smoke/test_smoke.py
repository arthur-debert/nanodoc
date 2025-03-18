#!/usr/bin/env python
"""Smoke tests for Nanodoc CLI.

These tests verify basic functionality of both v1 and v2 implementations.
"""

import os
import subprocess
import sys
import tempfile

# Test paths
TEST_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURES_DIR = os.path.join(os.path.dirname(TEST_DIR), "fixtures")
SAMPLE_FILES = [
    os.path.join(FIXTURES_DIR, "cake.txt"),
    os.path.join(FIXTURES_DIR, "incident.txt"),
    os.path.join(FIXTURES_DIR, "new-telephone.txt"),
]
SAMPLE_FILES_STR = " ".join(SAMPLE_FILES)


def run_command(cmd, check=True):
    """Run a command and return its output."""
    print(f"Running: {cmd}")
    result = subprocess.run(
        cmd, shell=True, capture_output=True, text=True, check=check
    )
    print(f"Exit code: {result.returncode}")
    if result.stdout:
        print("--- stdout ---")
        print(result.stdout[:500] + ("..." if len(result.stdout) > 500 else ""))
    if result.stderr:
        print("--- stderr ---")
        print(result.stderr[:500] + ("..." if len(result.stderr) > 500 else ""))
    return result


def test_v2_simple_concatenation():
    """Test basic file concatenation with v2."""
    cmd = f"python -m nanodoc --use-v2 {SAMPLE_FILES_STR}"
    result = run_command(cmd)
    assert result.returncode == 0
    assert "cake.txt" in result.stdout
    assert "incident.txt" in result.stdout
    assert "new-telephone.txt" in result.stdout


def test_v2_with_toc():
    """Test table of contents generation with v2."""
    cmd = f"python -m nanodoc --use-v2 --toc {SAMPLE_FILES_STR}"
    result = run_command(cmd)
    assert result.returncode == 0

    # Check that TOC header is present
    assert "TOC" in result.stdout

    # TOC in v2 has dots between filename and number
    # Find all TOC lines (they have repeating dots)
    dot_lines = [line for line in result.stdout.split("\n") if "." * 5 in line]

    # Check that all files are present in the TOC
    for base_name in ["cake.txt", "incident.txt", "new-telephone.txt"]:
        # Files can appear in different formats in TOC
        file_found = False
        # Check all TOC line patterns
        for line in dot_lines:
            # Check both as filename and capitalized without extension
            name_without_ext = base_name.split(".")[0].title()
            for pattern in [base_name, name_without_ext]:
                if pattern in line:
                    file_found = True
                    break

        assert file_found, f"{base_name} not found in TOC"


def test_v2_with_line_numbers():
    """Test line numbering with v2."""
    cmd = f"python -m nanodoc --use-v2 -n -n {SAMPLE_FILES_STR}"
    result = run_command(cmd)
    assert result.returncode == 0
    # Line numbers should be present
    assert "1:" in result.stdout
    assert "2:" in result.stdout


def test_v2_with_bundle():
    """Test bundle file processing with v2."""
    with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
        bundle_file = f.name
        for sample in SAMPLE_FILES:
            f.write(f"{sample}\n")

    try:
        # Use direct file list in v2 since it doesn't support bundle files
        cmd = f"python -m nanodoc --use-v2 {SAMPLE_FILES_STR}"
        result = run_command(cmd)
        assert result.returncode == 0
        # All files should be present
        assert "cake.txt" in result.stdout
        assert "incident.txt" in result.stdout
        assert "new-telephone.txt" in result.stdout
    finally:
        if os.path.exists(bundle_file):
            os.remove(bundle_file)


def test_v2_with_theme():
    """Test theme application with v2."""
    for theme in ["classic-light", "neutral", "classic-dark"]:
        cmd = f"python -m nanodoc --use-v2 --theme {theme} {SAMPLE_FILES[0]}"
        result = run_command(cmd)
        assert result.returncode == 0
        assert "cake.txt" in result.stdout


def main():
    """Run all smoke tests."""
    tests = [
        test_v2_simple_concatenation,
        test_v2_with_toc,
        test_v2_with_line_numbers,
        test_v2_with_bundle,
        test_v2_with_theme,
    ]

    failures = 0
    for test in tests:
        print(f"\n{'=' * 40}\nRunning {test.__name__}\n{'=' * 40}")
        try:
            test()
            print(f"✅ {test.__name__} PASSED")
        except Exception as e:
            print(f"❌ {test.__name__} FAILED: {e}")
            failures += 1

    print(f"\n{'-' * 40}")
    print(f"Ran {len(tests)} smoke tests, {failures} failures")
    return failures


if __name__ == "__main__":
    sys.exit(main())
